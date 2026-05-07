class Client {
  final String id;
  final String name;

  Client({required this.id, required this.name});

  // Aceasta este "poarta" prin care trece JSON-ul de la Go în Dart
  factory Client.fromJson(Map<String, dynamic> json) {
    return Client(
      id: (json['id'] ?? '').toString(),
      name: (json['name'] ?? 'Нет имени').toString(),
    );
  }
}