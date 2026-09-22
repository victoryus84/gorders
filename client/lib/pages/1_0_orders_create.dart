import 'package:flutter/material.dart';
import 'package:flutter/foundation.dart';
import '../models/client.dart';
import '../models/contract.dart';
import '../services/api_service.dart';
import '../controllers/order_controller.dart';
import '../core/logger.dart';

class OrdersCreatePage extends StatefulWidget {
  final String title;
  const OrdersCreatePage({super.key, required this.title});

  @override
  State<OrdersCreatePage> createState() => _OrdersCreatePageState();
}

class _OrdersCreatePageState extends State<OrdersCreatePage> {
  // Folosim late pentru a inițializa controller-ul o singură dată
  late final OrderCreateController _controller;
  final ApiService _api = ApiService();

  @override
  void initState() {
    super.initState();
    _controller = OrderCreateController();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text(widget.title), centerTitle: true),
      body: ListenableBuilder(
        listenable: _controller,
        builder: (context, _) {
          return GestureDetector(
            onTap: () => FocusScope.of(context).unfocus(),
            child: Column(
              // FĂRĂ SingleChildScrollView!
              children: [
                // 1. Zona de sus (Dropdown-uri și Categorii)
                Padding(
                  padding: const EdgeInsets.all(12.0),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.stretch,
                    children: [
                      _buildPaymentDropdown(),
                      const SizedBox(height: 10),
                      _buildClientAutocomplete(),
                      const SizedBox(height: 10),
                      _buildContractDropdown(
                        key: ValueKey(
                          'contract_${_controller.selectedClient?.id}',
                        ),
                      ),
                      const SizedBox(height: 15),
                      _buildCategoriesBar(), // Categoriile apar o singură dată aici!
                    ],
                  ),
                ),

                // 2. Lista de produse (Se va derula ultrarapid în spațiul rămas)
                Expanded(child: _buildProductsList()),

                // 3. Butonul de jos (Fixat deasupra tastaturii)
                Padding(
                  padding: const EdgeInsets.all(12.0),
                  child: _buildSubmitButton(),
                ),
              ],
            ),
          );
        },
      ),
    );
  }

  // --- COMPONENTE ---

  Widget _buildPaymentDropdown() {
    return DropdownButtonFormField<String>(
      initialValue: _controller
          .paymentType, // Folosim initialValue pentru consistență cu controller-ul
      decoration: const InputDecoration(
        labelText: "Тип оплаты",
        prefixIcon: Icon(Icons.payments_outlined),
        border: OutlineInputBorder(),
      ),
      items: const [
        DropdownMenuItem(value: "Нал", child: Text("Наличные (Нал)")),
        DropdownMenuItem(value: "Бнал", child: Text("Безналичные (Бнал)")),
      ],
      onChanged: (val) => _controller.setPaymentType(val),
    );
  }

  Widget _buildClientAutocomplete() {
    return Autocomplete<Client>(
      displayStringForOption: (Client c) => c.name,
      optionsBuilder: (TextEditingValue textValue) async {
        if (textValue.text.length < 3) return const Iterable<Client>.empty();

        final messenger = ScaffoldMessenger.of(context);
        final rezultate = await _api.searchClients(textValue.text);

        if (!mounted) return rezultate;

        if (kDebugMode && rezultate.isEmpty) {
          messenger.showSnackBar(
            SnackBar(
              content: Text("Ничего не найдено pentru '${textValue.text}'"),
              duration: const Duration(seconds: 1),
            ),
          );
        }
        return rezultate;
      },
      onSelected: (Client selection) {
        _controller.selectClient(selection);
        myLog("✅ Client selectat: ${selection.name}");
      },
      fieldViewBuilder: (context, textController, focusNode, onFieldSubmitted) {
        return TextFormField(
          controller: textController,
          focusNode: focusNode,
          decoration: InputDecoration(
            labelText: "Клиент",
            hintText: "Min. 3 litere...",
            prefixIcon: const Icon(Icons.person_search),
            suffixIcon: textController.text.isNotEmpty
                ? IconButton(
                    icon: const Icon(Icons.clear),
                    onPressed: () => textController.clear(),
                  )
                : null,
            border: const OutlineInputBorder(),
          ),
        );
      },
    );
  }

  Widget _buildContractDropdown({Key? key}) {
    return DropdownButtonFormField<Contract>(
      key: key,
      isExpanded: true,

      // REZOLVARE WARNING: Folosim initialValue (Fără erori galbene!)
      initialValue: _controller.selectedContract,

      decoration: InputDecoration(
        labelText: "Договор",
        border: const OutlineInputBorder(),
        contentPadding: const EdgeInsets.symmetric(
          horizontal: 12,
          vertical: 15,
        ),
        prefixIcon: const Icon(Icons.description_outlined),
        suffixIcon: _controller.isLoadingContracts
            ? const Padding(
                padding: EdgeInsets.all(12.0),
                child: CircularProgressIndicator(strokeWidth: 2),
              )
            : null,
      ),
      hint: Text(
        _controller.selectedClient == null
            ? "Сначала выберите клиента"
            : (_controller.availableContracts.isEmpty &&
                  !_controller.isLoadingContracts)
            ? "Нет договоров"
            : "Выберите договор",
        overflow: TextOverflow.ellipsis,
      ),
      items: _controller.availableContracts.map((c) {
        return DropdownMenuItem<Contract>(
          value: c,
          child: LayoutBuilder(
            builder: (context, constraints) {
              return SizedBox(
                width: constraints.maxWidth,
                child: Text(
                  c.name,
                  overflow: TextOverflow.ellipsis,
                  maxLines: 1,
                  softWrap: false,
                ),
              );
            },
          ),
        );
      }).toList(),
      onChanged: _controller.isLoadingContracts
          ? null
          : (val) => _controller.selectContract(val),
    );
  }

  Widget _buildCategoriesBar() {
    // Dacă nu avem produse încă (se încarcă), nu arătăm nici bara
    if (_controller.availableProducts.isEmpty) return const SizedBox.shrink();

    return SizedBox(
      height: 45, // Înălțime fixă pentru butoanele orizontale
      child: ListView.builder(
        scrollDirection: Axis.horizontal,
        itemCount: _controller.categories.length,
        itemBuilder: (context, index) {
          final category = _controller.categories[index];
          final isSelected = _controller.selectedCategory == category;

          return Padding(
            padding: const EdgeInsets.only(right: 8.0),
            child: ChoiceChip(
              label: Text(category),
              selected: isSelected,
              selectedColor: Colors.blue,
              labelStyle: TextStyle(
                color: isSelected ? Colors.white : Colors.black,
                fontWeight: isSelected ? FontWeight.bold : FontWeight.normal,
              ),
              onSelected: (bool selected) {
                if (selected) {
                  _controller.selectCategory(category);
                }
              },
            ),
          );
        },
      ),
    );
  }

  Widget _buildProductsList() {
    if (_controller.availableProducts.isEmpty) {
      return const Center(child: CircularProgressIndicator());
    }

    return ListView.builder(
      // Atenție: NU mai punem shrinkWrap și nici NeverScrollableScrollPhysics
      padding: const EdgeInsets.symmetric(horizontal: 12.0),
      itemCount: _controller.filteredProducts.length,
      itemBuilder: (context, index) {
        final product = _controller.filteredProducts[index];
        final currentQty = _controller.getQuantity(product.id);

        return Card(
          elevation: 1,
          margin: const EdgeInsets.symmetric(vertical: 4),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(10),
          ),
          child: Padding(
            padding: const EdgeInsets.all(10.0),
            child: Row(
              children: [
                Expanded(
                  flex: 3,
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        product.name,
                        style: const TextStyle(
                          fontWeight: FontWeight.w600,
                          fontSize: 14,
                        ),
                        maxLines: 2,
                        overflow: TextOverflow.ellipsis,
                      ),
                      const SizedBox(height: 4),
                      Text(
                        "${product.price} MDL",
                        style: TextStyle(color: Colors.grey[700], fontSize: 13),
                      ),
                    ],
                  ),
                ),
                Expanded(
                  flex: 2,
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.end,
                    children: [
                      IconButton(
                        icon: const Icon(Icons.remove_circle_outline),
                        color: currentQty > 0 ? Colors.red : Colors.grey,
                        onPressed: currentQty > 0
                            ? () => _controller.updateQuantity(
                                product.id,
                                currentQty - 1,
                              )
                            : null,
                      ),
                      SizedBox(
                        width: 25,
                        child: Text(
                          "$currentQty",
                          textAlign: TextAlign.center,
                          style: const TextStyle(
                            fontSize: 16,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                      ),
                      IconButton(
                        icon: const Icon(Icons.add_circle_outline),
                        color: Colors.green,
                        onPressed: () => _controller.updateQuantity(
                          product.id,
                          currentQty + 1,
                        ),
                      ),
                    ],
                  ),
                ),
              ],
            ),
          ),
        );
      },
    );
  }

  Widget _buildSubmitButton() {
    final bool isReady = _controller.isValid && !_controller.isSubmitting;

    return ElevatedButton.icon(
      onPressed: isReady ? () => _handleSave() : null,
      icon: _controller.isSubmitting
          ? const SizedBox(
              width: 20,
              height: 20,
              child: CircularProgressIndicator(
                color: Colors.white,
                strokeWidth: 2,
              ),
            )
          : const Icon(Icons.save),
      label: Text(
        _controller.isSubmitting ? "ОТПРАВКА..." : "СОЗДАТЬ ЗАКАЗ",
        style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 16),
      ),
      style: ElevatedButton.styleFrom(
        minimumSize: const Size.fromHeight(55),
        backgroundColor: isReady ? Colors.blue : Colors.grey[300],
        foregroundColor: Colors.white,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(10)),
        elevation: isReady ? 3 : 0,
      ),
    );
  }

  Future<void> _handleSave() async {
    setState(() => _controller.isSubmitting = true);

    final jsonData = _controller.buildOrderJson();
    myLog("🚀 Trimitere JSON către Go:\n$jsonData");

    // Apelează serverul
    final success = await _api.submitOrder(jsonData);

    setState(() => _controller.isSubmitting = false);

    if (mounted) {
      if (success) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text("✅ Заказ успешно создан!"),
            backgroundColor: Colors.green,
          ),
        );
        Navigator.pop(context); // Opțional: te întoarce la pagina anterioară
      } else {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text("❌ Ошибка сервера!"),
            backgroundColor: Colors.red,
          ),
        );
      }
    }
  }
}
