class Contract {
  final String id;
  final String name;

  Contract({required this.id, required this.name});

  factory Contract.fromJson(Map<String, dynamic> json) {
    return Contract(
      // ID-ul (încercăm și cu mic, și cu mare)
      id: (json['id'] ?? json['ID'] ?? '').toString(),

      // NUMELE: Încercăm toate variantele posibile de la Go/1C
      name:
          (json['name'] ?? // varianta json mica
                  json['Name'] ?? // varianta Go standard
                  json['number'] ?? // poate e numarul contractului
                  json['Number'] ?? // varianta 1C/Go mare
                  'Договор без названия') // fallback
              .toString(),
    );
  }
}
