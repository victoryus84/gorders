// Список заказов

import 'package:flutter/material.dart';
import 'package:frontend/widgets/back_menu.dart';

class OrdersListPage extends StatelessWidget {
  final String title;
  const OrdersListPage({super.key, required this.title});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text(title), actions: [const BackToMenuButton()]),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Form(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              //Text("Страница: $title", style: const TextStyle(fontSize: 18)),
              const SizedBox(height: 20),
              TextFormField(
                decoration: const InputDecoration(labelText: "Enter something"),
              ),
              const SizedBox(height: 20),
              ElevatedButton(
                onPressed: () {
                  ScaffoldMessenger.of(context).showSnackBar(
                    const SnackBar(content: Text("Form submitted!")),
                  );
                },
                child: const Text("Submit"),
              ),
              const SizedBox(height: 20),
            ],
          ),
        ),
      ),
    );
  }
}
