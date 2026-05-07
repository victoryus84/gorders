import 'package:flutter/material.dart';
import '../models/client.dart';
import '../models/contract.dart';
import '../services/api_service.dart';
import '../core/logger.dart';

class OrderCreateController extends ChangeNotifier {
  final ApiService _api = ApiService();

  Client? selectedClient;
  Contract? selectedContract;
  String? paymentType;

  List<Contract> availableContracts = [];
  bool isLoadingContracts = false;

  // Verifică dacă butonul de salvare poate fi apăsat
  bool get isValid =>
      selectedClient != null && selectedContract != null && paymentType != null;

  void setPaymentType(String? type) {
    paymentType = type;
    notifyListeners();
  }

  void selectContract(Contract? contract) {
    selectedContract = contract;
    notifyListeners();
  }

  // LOGICA DE DEPENDENȚĂ: Client -> Contract
  Future<void> selectClient(Client? client) async {
    if (client == null) {
      // DACĂ ȘTERGEM CLIENTUL:
      selectedClient = null;
      selectedContract = null;
      availableContracts = [];
      isLoadingContracts = false; // Oprim rotița de încărcare
      notifyListeners();
      return; // Ne oprim aici, nu mai chemăm API-ul
    }

    // DACĂ AVEM UN CLIENT NOU:
    selectedClient = client;
    selectedContract = null;
    availableContracts = [];
    isLoadingContracts = true;
    notifyListeners();

    try {
      availableContracts = await _api.fetchContracts(client.id);
    } catch (e) {
      myLog("❌ Eroare la încărcare contracte: $e");
    } finally {
      isLoadingContracts = false;
      notifyListeners();
    }
  }
}
