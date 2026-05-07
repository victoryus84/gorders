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
            // Închide tastatura când dai click pe fundal
            onTap: () => FocusScope.of(context).unfocus(),
            child: SingleChildScrollView(
              padding: const EdgeInsets.all(16.0),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment
                    .stretch, // Face butoanele lățime completă
                children: [
                  _buildPaymentDropdown(),
                  const SizedBox(height: 20),
                  _buildClientAutocomplete(),
                  const SizedBox(height: 20),

                  // Contract Dropdown cu resetare automată la schimbarea clientului
                  _buildContractDropdown(
                    key: ValueKey(
                      'contract_${_controller.selectedClient?.id ?? 'none'}',
                    ),
                  ),

                  const SizedBox(height: 40),
                  _buildSubmitButton(),
                ],
              ),
            ),
          );
        },
      ),
    );
  }

  // --- COMPONENTE ---

  Widget _buildPaymentDropdown() {
    return DropdownButtonFormField<String>(
      initialValue: _controller.paymentType, // Folosim initialValue pentru consistență cu controller-ul
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
            : "Выберите договор",
        overflow: TextOverflow.ellipsis,
      ),
      // Mapăm lista de contracte cu protecție la Overflow
      items: _controller.availableContracts.map((c) {
        return DropdownMenuItem<Contract>(
          value: c,
          child: LayoutBuilder(
            // Folosim LayoutBuilder pentru a afla lățimea exactă disponibilă
            builder: (context, constraints) {
              return SizedBox(
                width: constraints
                    .maxWidth, // Forțează textul să știe cât spațiu are
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

  Widget _buildSubmitButton() {
    final bool isReady = _controller.isValid;
    return ElevatedButton.icon(
      onPressed: isReady ? () => _handleSave() : null,
      icon: const Icon(Icons.save),
      label: const Text(
        "СОЗДАТЬ ЗАКАЗ",
        style: TextStyle(fontWeight: FontWeight.bold),
      ),
      style: ElevatedButton.styleFrom(
        minimumSize: const Size.fromHeight(55),
        backgroundColor: isReady ? Colors.green : Colors.grey[300],
        foregroundColor: Colors.white,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
      ),
    );
  }

  void _handleSave() {
    myLog(
      "🚀 Salvare comandă: Client: ${_controller.selectedClient?.name}, Contract: ${_controller.selectedContract?.name}",
    );
    // Aici vine apelul tău de POST către Go
  }
}
