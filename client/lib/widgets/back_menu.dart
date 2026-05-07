import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

/// Универсальная кнопка возврата в главное меню
class BackToMenuButton extends StatelessWidget {
  const BackToMenuButton({super.key});

  @override
  Widget build(BuildContext context) {
    return ElevatedButton.icon(
      onPressed: () {
        context.go('/0_1_menu'); // всегда ведёт в меню
      },
      icon: const Icon(Icons.home),
      label: const Text("Menu"),
      style: ElevatedButton.styleFrom(
        backgroundColor: Colors.lightGreen,
        padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 12),
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(10)),
      ),
    );
  }
}
