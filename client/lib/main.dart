import 'dart:io';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:google_fonts/google_fonts.dart';

import '/services/auth_service.dart';
import 'pages/0_0_login.dart';
import 'pages/0_1_menu.dart';
import 'pages/1_0_orders_create.dart';
import 'pages/1_1_orders_list.dart';
import 'pages/2_0_equipment.dart';
import 'pages/3_0_calculations.dart';
import 'pages/4_0_shipment.dart';

Future<void> main() async {
  // Инициализируем AuthService до запуска приложения
  WidgetsFlutterBinding.ensureInitialized();
  await AuthService.init();

  HttpOverrides.global = MyHttpOverrides();
  
  runApp(const MyApp());
}

class MyHttpOverrides extends HttpOverrides {
  @override
  HttpClient createHttpClient(SecurityContext? context) {
    return super.createHttpClient(context)
      ..badCertificateCallback = (X509Certificate cert, String host, int port) {
        // Разрешаем самоподписанные сертификаты для вашего домена и IP
        return host == 'servidar.work.gd' || host == '217.26.172.96';
      };
  }
}

final _router = GoRouter(
  initialLocation: '/0_0_login',
  routes: [
    GoRoute(path: '/0_0_login', builder: (context, state) => const LoginPage()),
    GoRoute(path: '/0_1_menu', builder: (context, state) => const MenuPage()),
    GoRoute(path: '/1_0_orders_create', builder: (context, state) => OrdersCreatePage(title: "Создать заявку"),),
    GoRoute(path: '/1_1_orders_list', builder: (context, state) => const OrdersListPage(title: "Список заявок"),),
    GoRoute(path: '/2_0_equipment', builder: (context, state) => const EquipmentPage(title: "Оборудование"),),
    GoRoute(path: '/3_0_calculations', builder: (context, state) => const CalculationsPage(title: "Взаиморасчеты"),),
    GoRoute(path: '/4_0_shipment', builder: (context, state) => const ShipmentPage(title: "Отгрузка 4 мес."),),
],

  /// 🔹 Синхронный redirect (используем кешированный токен)
  redirect: (context, state) {
    final token = AuthService.token; // Синхронный доступ к токену
    final loggingIn = state.matchedLocation == '/0_0_login';

    // Если нет токена и не на странице логина - на логин
    if (token == null && !loggingIn) {
      return '/0_0_login';
    }

    // Если есть токен и на странице логина - в меню
    if (token != null && loggingIn) {
      return '/0_1_menu';
    }

    // Во всех остальных случаях - не перенаправлять
    return null;
  },
);

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  // This widget is the root of your application.
  @override
  Widget build(BuildContext context) {
    return MaterialApp.router(
      title: 'Orders App',
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.lightGreen),
        textTheme: GoogleFonts.robotoTextTheme(
          Theme.of(context).textTheme,
        ),
      ),
      routerConfig: _router,
    );
  }
}
