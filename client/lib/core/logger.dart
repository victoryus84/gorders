import 'dart:developer' as dev;
import 'package:flutter/foundation.dart'; // Obligatoriu pentru kDebugMode

void myLog(String message, {Object? error, String name = 'app.orders'}) {
  // Verifică dacă ești în modul Debug (când rulezi de pe laptop)
  if (kDebugMode) {
    dev.log(message, error: error, name: name, time: DateTime.now());
  }
}
