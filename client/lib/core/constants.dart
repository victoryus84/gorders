class AppConfig {
  // Adresa ta de bază
  static const String baseUrl = "https://servidar.work.gd";

  // Poți pune aici și alte rute ca să fie totul într-un loc
  static const String clientsEndpoint = "$baseUrl/api/v1/clients/search";
  static const String clientsSearchEndpoint = "$baseUrl/api/v1/clients/search";
  static const String ordersEndpoint = "$baseUrl/api/v1/orders";
  static const String contractsEndpointByClient = "$baseUrl/api/v1/contracts/client";
  // Timeout-uri sau alte setări globale
  // Timeout-ul implicit pentru cererile API
  static const int apiTimeout = 5000; // milisecunde
}
