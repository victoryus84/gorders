import 'package:flutter/material.dart';
import '../models/client.dart';
import '../models/contract.dart';
import '../models/product.dart';
import '../services/api_service.dart';
import '../core/logger.dart';

class OrderCreateController extends ChangeNotifier {
  final ApiService _api = ApiService();

  // --- Stări pentru Formular ---
  String paymentType = "Нал";
  Client? selectedClient;
  Contract? selectedContract;

  List<Contract> availableContracts = [];
  bool isLoadingContracts = false;

  // --- Stări pentru Produse ---
  List<Product> availableProducts = [];

  String selectedCategory = 'Все'; // Categoria selectată implicit (Toate)

  // 1. Extragem automat categoriile unice din lista de produse
  List<String> get categories {
    final cats = availableProducts.map((p) => p.category).toSet().toList();
    cats.insert(0, 'Все'); // Adăugăm opțiunea "Toate" la început
    return cats;
  }

  // 2. Returnăm doar produsele din categoria selectată
  List<Product> get filteredProducts {
    if (selectedCategory == 'Все') return availableProducts;
    return availableProducts
        .where((p) => p.category == selectedCategory)
        .toList();
  }

  // 3. Funcția care schimbă categoria când apeși pe un buton
  void selectCategory(String category) {
    selectedCategory = category;
    notifyListeners(); // Redesenează UI-ul instant
  }

  final Map<String, int> _quantities = {}; // Salvează cantitatea per produs
  bool isSubmitting = false;

  // Constructor: Poți încărca produsele automat când se deschide pagina
  OrderCreateController() {
    _loadProducts();
  }

  Future<void> _loadProducts() async {
    // Aici chemi API-ul tău pentru a aduce lista de produse
    // availableProducts = await _api.fetchProducts();
    availableProducts = await _api.fetchProducts();
    notifyListeners();
  }

  // --- Metode UI ---
  void setPaymentType(String? val) {
    if (val != null) {
      paymentType = val;
      notifyListeners();
    }
  }

  Future<void> selectClient(Client? client) async {
    if (client == null) {
      selectedClient = null;
      selectedContract = null;
      availableContracts = [];
      isLoadingContracts = false;
      notifyListeners();
      return;
    }

    selectedClient = client;
    selectedContract = null;
    availableContracts = [];
    isLoadingContracts = true;
    notifyListeners();

    try {
      availableContracts = await _api.fetchContracts(client.id);
      myLog(
        "📦 Contracte aduse pentru ${client.name}: ${availableContracts.length}",
      );
    } catch (e) {
      myLog("❌ Eroare la contracte: $e");
    } finally {
      isLoadingContracts = false;
      notifyListeners();
    }
  }

  void selectContract(Contract? contract) {
    selectedContract = contract;
    notifyListeners();
  }

  // --- Metode Produse ---
  int getQuantity(String productId) => _quantities[productId] ?? 0;

  void updateQuantity(String productId, int quantity) {
    if (quantity <= 0) {
      _quantities.remove(productId);
    } else {
      _quantities[productId] = quantity;
    }
    notifyListeners();
  }

  // --- Validare și Trimitere ---

  // E valid doar dacă avem Client și MĂCAR un produs cu cantitate > 0
  bool get isValid =>
      selectedClient != null && _quantities.values.any((q) => q > 0);

  Map<String, dynamic> buildOrderJson() {
    final items = _quantities.entries
        .map((e) => {'product_id': e.key, 'quantity': e.value})
        .toList();

    return {
      'client_id': selectedClient?.id,
      'contract_id': selectedContract
          ?.id, // Poate fi null, backend-ul Go trebuie să știe asta
      'payment_type': paymentType,
      'items': items,
    };
  }
}
