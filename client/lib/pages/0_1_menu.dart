import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import '/services/auth_service.dart'; // Импортируем AuthService

/// Main menu form with 5 buttons
class MenuPage extends StatelessWidget {
  const MenuPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text("Основное Menu")),
      backgroundColor: Colors.grey[200],
      body: Center(
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            const SizedBox(height: 20),
            ElevatedButton(
              onPressed: () => context.go('/1_0_orders_create'),
              style: ElevatedButton.styleFrom(
                padding: const EdgeInsets.symmetric(vertical: 16),
              ),
              child: const Text("Создать заявку"),
            ),
            const SizedBox(height: 12),
            ElevatedButton(
              onPressed: () => context.go('/1_1_orders_list'),
              style: ElevatedButton.styleFrom(
                padding: const EdgeInsets.symmetric(vertical: 16),
              ),
              child: const Text("Список заявок"),
            ),
            const SizedBox(height: 12),
            ElevatedButton(
              onPressed: () => context.go('/2_0_equipment'),
              style: ElevatedButton.styleFrom(
                padding: const EdgeInsets.symmetric(vertical: 16),
              ),
              child: const Text("Оборудование"),
            ),
            const SizedBox(height: 12),
            ElevatedButton(
              onPressed: () => context.go('/3_0_calculations'),
              style: ElevatedButton.styleFrom(
                padding: const EdgeInsets.symmetric(vertical: 16),
              ),
              child: const Text("Взаиморасчеты"),
            ),
            const SizedBox(height: 12),
            ElevatedButton(
              onPressed: () => context.go('/4_0_shipment'),
              style: ElevatedButton.styleFrom(
                padding: const EdgeInsets.symmetric(vertical: 16),
              ),
              child: const Text("Отгрузка 4 мес."),
            ),
            const SizedBox(height: 20),
            ElevatedButton(
              onPressed: () async {
                await AuthService.clearToken();
                // ignore: use_build_context_synchronously
                GoRouter.of(context).go('/0_0_login');
              },
              child: const Text('Выйти'),
            )
          ],
        ),
      ),
    );
  }
}
