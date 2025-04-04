= SPEC-1: Информационная система управления Доставкой
:sectnums:
:toc:

== Background

Современные сервисы доставки требуют эффективного управления заказами, оплатами, меню и пользователями. Информационная система должна обеспечивать взаимодействие клиентов, официантов, поваров и администраторов через удобный интерфейс. Важными аспектами являются мониторинг, логирование и аналитика данных.

== Requirements

=== Must-have (Обязательные)

- Регистрация и аутентификация пользователей (клиенты, администраторы, официанты, повара).
- Управление заказами (создание, обновление, статусы, уведомления).
- Интеграция с платежной системой для онлайн-оплат (заглушка).
- Управление меню (CRUD-операции с блюдами, категориями, ценами).
- Система отчетности и аналитики (реплики баз данных для BI-инструмента Metabase).
- Логирование событий и мониторинг работы системы.

=== Should-have (Желательные)

- Поддержка программ лояльности и скидок.
- Интеграция с внешними сервисами доставки.

=== Could-have (Дополнительные)

- Автоматизация обработки заказов с помощью ML-модели (например, предсказание популярности блюд). (Не важно, будет алгоритмом, без ML)

=== Won't-have (Не планируется на текущий этап)

- Поддержка нескольких ресторанов в одной системе.
- Голосовое управление заказами.

== Method

=== Архитектурный стиль

- Микросервисная архитектура.
- API Gateway для управления запросами.
- Асинхронное взаимодействие (Kafka / RabbitMQ для событий).
- Репликация БД для аналитики (Metabase).
- Паттерн SAGA для контроля за операциями

=== Технологический стек

- **Бэкенд:** Go (Gin, gRPC, pgxPool)
- **Фронтенд:** React (Next.js / Vite)
- **База данных:** PostgreSQL (отдельная БД на сервис), S3 для хранения фото блюд. Мб кликхаус на нагруженные бд вместо PostgreSQL
- **Сообщения:** Kafka
- **Мониторинг:** Prometheus + Grafana
- **Логирование:** Loki + Grafana
- **Развертывание** Docker + IaC(Ansible, Terraform)

=== Основные сервисы

0. **Notify** - Сервис для отправки уведомлений пользователям(будет предствалять из себя заглушку с тг ботом)
1. **AuthService** – Аутентификация (JWT). Принимает логин, пароль для base_auth, либо интеграция с yandex Oauth2 для пользователя, отдает JWT токен, с которым пользователь может пользоваться другими сервисами
   Ответственность:
   Аутентификация и авторизация пользователей/персонала.

Основные функции:

Генерация JWT-токенов (срок жизни: 24 часа).

Интеграция с Yandex OAuth2 для клиентов.

Проверка ролей (customer, admin, chef, waiter).

Хеширование паролей (bcrypt для пользователей, Argon2 для персонала).

API-эндпоинты (gRPC):

rpc Login(LoginRequest) returns (AuthResponse) – вход по логину/паролю.

rpc OAuthYandex(OAuthRequest) returns (AuthResponse) – вход через Yandex.

rpc ValidateToken(TokenRequest) returns (RoleResponse) – проверка токена.

2. **UserService** – Управление пользователями. Создание, хранение, удаление.
   Ответственность:
   Управление данными клиентов и персонала.

Основные функции:

CRUD для клиентов (users таблица).

CRUD для персонала (staff таблица).

Назначение ролей через user_roles и staff_roles.

Валидация email/телефона (уникальность).

API-эндпоинты (gRPC):

rpc CreateUser(UserRequest) returns (UserResponse)

rpc UpdateUser(UserUpdate) returns (Empty)

rpc GetStaff(StaffID) returns (StaffResponse)

Важно:

Сервис не хранит пароли – только user_id для связи с AuthService.

3. **MenuService** – Управление меню.

Ответственность:
Управление меню, блюдами, категориями и акциями.

Основные функции:

CRUD для блюд, категорий, меню.

Расчет цены с учетом акций (promotions).

Генерация pre-signed URLs для загрузки фото в S3.

Кэширование меню в Redis (TTL: 1 час).

API-эндпоинты (gRPC):

rpc GetActiveMenu(Empty) returns (MenuResponse)

rpc UpdateDish(DishUpdate) returns (Empty)

rpc ApplyPromotion(PromotionRequest) returns (Empty)

4. **OrderService** – Управление заказами.

Ответственность:
Оформление и отслеживание заказов.

Основные функции:

Создание заказа с валидацией доступности блюд.

Обновление статусов: created → confirmed → ready → delivered.

Интеграция с PaymentService через SAGA (Kafka).

Расчет стоимости заказа с учетом скидок.

API-эндпоинты (gRPC):

rpc CreateOrder(OrderRequest) returns (OrderResponse)

rpc CancelOrder(OrderID) returns (Empty)

rpc GetOrderStatus(OrderID) returns (StatusResponse)

События Kafka:

order_created – запуск процесса оплаты.

order_failed – отмена из-за ошибки оплаты.

При отмене платежа OrderService должен:

Обновить статус заказа на canceled.

Отправить событие order_failed в Kafka для уведомления клиента.

5. **PaymentService** – Обработка платежей (заглушка).

Ответственность:
Обработка платежей (заглушка с логикой "успех/ошибка").

Основные функции:

Имитация платежных шлюзов (Stripe/YooKassa).

Генерация случайных транзакционных ID.

Интеграция с OrderService через SAGA.

Возврат средств при отмене заказа.

API-эндпоинты (gRPC):

rpc ProcessPayment(PaymentRequest) returns (PaymentResponse)

rpc Refund(OrderID) returns (Empty)

6. **Gateway** - начальная точка входа в приложение, роутинг, конвертация gRPC ответов от сервисов в Http ответы и обратно

Ответственность:
Единая точка входа для фронтенда.

Основные функции:

Маршрутизация HTTP-запросов к gRPC-сервисам.

Конвертация gRPC ↔ JSON.

Аутентификация через JWT (вызов AuthService).

Rate limiting (100 RPM на пользователя).

Маршруты (HTTP):

POST /api/orders → OrderService.CreateOrder

GET /api/menu → MenuService.GetActiveMenu

POST /api/auth/login → AuthService.Login

Промежуточное ПО:

Проверка JWT для /api/\* (кроме /auth/login).

Кэширование GET-запросов к /api/menu.

7. **UserFrontend** - Web приложения для пользователя, заказы, просмотр меню, регистрация, просмотр статуса заказа, оплата()
   Ответственность:
   Web-интерфейс для клиентов.

Стек: React + Vite + Tailwind CSS.

Страницы:

/menu – просмотр блюд с фильтрами по категориям.

/cart – оформление заказа (выбор типа доставки).

/orders – история заказов со статусами.

/payment – форма оплаты (карта/наличные).

Интеграция:

Запросы к Gateway через Axios.

WebSocket для实时 обновления статусов заказов.

8. **PersonalFrontend** - Web приложение для работников, оформление заказа, подтверждение, передача на кухню, отметка о готовности,

Ответственность:
Web-интерфейс для сотрудников.

Стек: Next.js + Material UI.

Страницы:

/dashboard – список активных заказов (фильтры по статусам).

/kitchen – управление готовностью блюд (кнопка "Готово").

/analytics – базовые отчеты (Metabase iframe).

Фичи:

Drag-and-drop для изменения порядка заказов.

Уведомления через Toast при новых заказах.

=== Структура APIS

**AuthService:**

```proto
syntax = "proto3";

package auth.v1;
option go_package = "github.com/netscrawler/Restaurant_is/gen/go/auth/v1;authv1";

import "google/protobuf/empty.proto";

service AuthService {
  // Базовая аутентификация (логин/пароль)
  rpc Login(LoginRequest) returns (LoginResponse);

  // OAuth2 через Yandex
  rpc LoginYandex(OAuthYandexRequest) returns (LoginResponse);

  // Валидация и декодирование JWT токена
  rpc Validate(ValidateRequest) returns (ValidateResponse);

  // Обновление токена (не реализовано в первой версии)
  rpc Refresh(RefreshRequest) returns (LoginResponse);
}

// Роли пользователей
enum Role {
  ROLE_UNSPECIFIED = 0;
  ROLE_CUSTOMER = 1;    // Обычный клиент
  ROLE_WAITER = 2;      // Официант
  ROLE_CHEF = 3;        // Повар
  ROLE_ADMIN = 4;       // Администратор
}

// Тип учетной записи
enum AccountType {
  ACCOUNT_TYPE_UNSPECIFIED = 0;
  ACCOUNT_TYPE_USER = 1;     // Клиент
  ACCOUNT_TYPE_STAFF = 2;    // Сотрудник
}

message LoginRequest {
  oneof identifier {
    string email = 1;        // Для пользователей
    string work_email = 2;   // Для сотрудников
  }
  string password = 3;
}

message OAuthYandexRequest {
  string code = 1;           // Код авторизации от Yandex OAuth
  string redirect_uri = 2;   // URI перенаправления
}

message LoginResponse {
  string access_token = 1;   // JWT токен
  int64 expires_in = 2;      // Время жизни токена в секундах
  User user = 3;
}

message ValidateRequest {
  string access_token = 1;
}

message ValidateResponse {
  bool valid = 1;
  User user = 2;
}

message RefreshRequest {
  string refresh_token = 1;
}

// Общая информация о пользователе
message User {
  int64 id = 1;
  AccountType account_type = 2;
  repeated Role roles = 3;
  UserMetadata metadata = 4;
}

// Дополнительные метаданные пользователя
message UserMetadata {
  string email = 1;
  string phone = 2;
  string full_name = 3;

  // Только для сотрудников
  string position = 10;
  string hire_date = 11;
}
```

**UserService:**

```proto
syntax = "proto3";

package user.v1;
option go_package = "github.com/your-project/gen/go/user/v1;userv1";

import "google/protobuf/empty.proto";
import "google/protobuf/timestamp.proto";

service UserService {
  // Пользователи (клиенты)
  rpc CreateUser(CreateUserRequest) returns (UserResponse);
  rpc GetUser(GetUserRequest) returns (UserResponse);
  rpc UpdateUser(UpdateUserRequest) returns (UserResponse);
  rpc DeleteUser(DeleteUserRequest) returns (google.protobuf.Empty);
  rpc ListUsers(ListUsersRequest) returns (ListUsersResponse);

  // Сотрудники
  rpc CreateStaff(CreateStaffRequest) returns (StaffResponse);
  rpc UpdateStaff(UpdateStaffRequest) returns (StaffResponse);
  rpc ListStaff(ListStaffRequest) returns (ListStaffResponse);

  // Роли
  rpc AssignRole(AssignRoleRequest) returns (google.protobuf.Empty);
  rpc RevokeRole(RevokeRoleRequest) returns (google.protobuf.Empty);
}

// Основные сообщения

message User {
  int64 id = 1;
  string email = 2;
  string phone = 3;
  string full_name = 4;
  bool is_active = 5;
  google.protobuf.Timestamp created_at = 6;
  google.protobuf.Timestamp updated_at = 7;
  repeated string roles = 8;
}

message Staff {
  int64 id = 1;
  string work_email = 2;
  string work_phone = 3;
  string full_name = 4;
  string position = 5;
  bool is_active = 6;
  google.protobuf.Timestamp hire_date = 7;
  repeated string roles = 8;
}

// Запросы/Ответы

message CreateUserRequest {
  string email = 1;
  string phone = 2;
  string full_name = 3;
  string password = 4; // Хешируется на сервере
}

message GetUserRequest {
  oneof identifier {
    int64 id = 1;
    string email = 2;
    string phone = 3;
  }
}

message UpdateUserRequest {
  int64 id = 1;
  optional string email = 2;
  optional string phone = 3;
  optional string full_name = 4;
  optional bool is_active = 5;
}

message DeleteUserRequest {
  int64 id = 1;
}

message ListUsersRequest {
  optional bool only_active = 1;
  int32 page = 2;
  int32 page_size = 3;
}

message ListUsersResponse {
  repeated User users = 1;
  int32 total_count = 2;
}

message CreateStaffRequest {
  string work_email = 1;
  string work_phone = 2;
  string full_name = 3;
  string position = 4;
  string password = 5;
}

message UpdateStaffRequest {
  int64 id = 1;
  optional string work_phone = 2;
  optional string position = 3;
  optional bool is_active = 4;
}

message ListStaffRequest {
  optional bool only_active = 1;
  int32 page = 2;
  int32 page_size = 3;
}

message ListStaffResponse {
  repeated Staff staff = 1;
  int32 total_count = 2;
}

message AssignRoleRequest {
  int64 user_id = 1;
  string role = 2;
  bool is_staff = 3; // true для сотрудников
}

message RevokeRoleRequest {
  int64 user_id = 1;
  string role = 2;
  bool is_staff = 3;
}

// Ответы

message UserResponse {
  User user = 1;
}

message StaffResponse {
  Staff staff = 1;
}
```

**MenuService:**

```proto
syntax = "proto3";

package menu.v1;
option go_package = "github.com/netscrawler/Restaurant_is/gen/go/menu/v1;menuv1";

import "google/protobuf/empty.proto";
import "google/protobuf/timestamp.proto";

service MenuService {
  // Категории
  rpc CreateCategory(CategoryRequest) returns (CategoryResponse);
  rpc UpdateCategory(UpdateCategoryRequest) returns (CategoryResponse);
  rpc ListCategories(ListCategoriesRequest) returns (ListCategoriesResponse);
  rpc DeleteCategory(DeleteCategoryRequest) returns (google.protobuf.Empty);

  // Блюда
  rpc CreateDish(DishRequest) returns (DishResponse);
  rpc UpdateDish(UpdateDishRequest) returns (DishResponse);
  rpc GetDish(GetDishRequest) returns (DishResponse);
  rpc ListDishes(ListDishesRequest) returns (ListDishesResponse);

  // Меню
  rpc CreateMenu(MenuRequest) returns (MenuResponse);
  rpc GetActiveMenu(google.protobuf.Empty) returns (MenuResponse);
  rpc UpdateMenu(UpdateMenuRequest) returns (MenuResponse);

  // Акции
  rpc CreatePromotion(PromotionRequest) returns (PromotionResponse);
  rpc UpdatePromotion(UpdatePromotionRequest) returns (PromotionResponse);
  rpc ListActivePromotions(google.protobuf.Empty) returns (ListPromotionsResponse);

  // Изображения
  rpc GenerateUploadUrl(ImageRequest) returns (ImageResponse);
}

// Основные сообщения

message Category {
  int64 id = 1;
  string name = 2;
  string description = 3;
  int32 display_order = 4;
  bool is_active = 5;
}

message Dish {
  int64 id = 1;
  string name = 2;
  string description = 3;
  double price = 4;
  int64 category_id = 5;
  int32 cooking_time_min = 6;
  string image_url = 7;
  bool is_available = 8;
  int32 calories = 9;
  google.protobuf.Timestamp created_at = 10;
  google.protobuf.Timestamp updated_at = 11;
}

message Menu {
  int64 id = 1;
  string name = 2;
  string description = 3;
  google.protobuf.Timestamp valid_from = 4;
  google.protobuf.Timestamp valid_to = 5;
  bool is_active = 6;
  repeated MenuCategory categories = 7;
  repeated int64 dish_ids = 8; // Блюда вне категорий
}

message MenuCategory {
  int64 category_id = 1;
  int32 display_order = 2;
}

message Promotion {
  int64 id = 1;
  int64 dish_id = 2;
  int32 discount_percent = 3; // 1-100%
  google.protobuf.Timestamp start_date = 4;
  google.protobuf.Timestamp end_date = 5;
}

// Запросы/Ответы

message CategoryRequest {
  string name = 1;
  string description = 2;
  int32 display_order = 3;
}

message UpdateCategoryRequest {
  int64 id = 1;
  optional string name = 2;
  optional string description = 3;
  optional int32 display_order = 4;
  optional bool is_active = 5;
}

message ListCategoriesRequest {
  bool only_active = 1;
  int32 page = 2;
  int32 page_size = 3;
}

message ListCategoriesResponse {
  repeated Category categories = 1;
  int32 total_count = 2;
}

message DeleteCategoryRequest {
  int64 id = 1;
  bool force = 2; // Принудительное удаление с блюдами
}

message DishRequest {
  string name = 1;
  string description = 2;
  double price = 3;
  int64 category_id = 4;
  int32 cooking_time_min = 5;
  optional string image_url = 6;
  bool is_available = 7;
  optional int32 calories = 8;
}

message UpdateDishRequest {
  int64 id = 1;
  optional string name = 2;
  optional string description = 3;
  optional double price = 4;
  optional int64 category_id = 5;
  optional int32 cooking_time_min = 6;
  optional string image_url = 7;
  optional bool is_available = 8;
  optional int32 calories = 9;
}

message ListDishesRequest {
  optional int64 category_id = 1;
  bool only_available = 2;
  int32 page = 3;
  int32 page_size = 4;
}

message ListDishesResponse {
  repeated Dish dishes = 1;
  int32 total_count = 2;
}

message MenuRequest {
  string name = 1;
  string description = 2;
  google.protobuf.Timestamp valid_from = 3;
  google.protobuf.Timestamp valid_to = 4;
  repeated MenuCategory categories = 5;
  repeated int64 dish_ids = 6;
}

message UpdateMenuRequest {
  int64 id = 1;
  optional string name = 2;
  optional string description = 3;
  optional google.protobuf.Timestamp valid_from = 4;
  optional google.protobuf.Timestamp valid_to = 5;
  optional bool is_active = 6;
  repeated MenuCategory categories = 7;
  repeated int64 dish_ids = 8;
}

message PromotionRequest {
  int64 dish_id = 1;
  int32 discount_percent = 2;
  google.protobuf.Timestamp start_date = 3;
  google.protobuf.Timestamp end_date = 4;
}

message ListPromotionsResponse {
  repeated Promotion promotions = 1;
}

message ImageRequest {
  string filename = 1; // "salad.jpg"
  string content_type = 2; // "image/jpeg"
}

message ImageResponse {
  string url = 1; // Pre-signed S3 URL
  string object_key = 2; // Для сохранения в БД
}

// Общие ответы
message CategoryResponse {
  Category category = 1;
}

message DishResponse {
  Dish dish = 1;
}

message MenuResponse {
  Menu menu = 1;
}

message PromotionResponse {
  Promotion promotion = 1;
}

```

**OrderService:**

```proto
syntax = "proto3";

package order.v1;
option go_package = "github.com/netscrawler/Restaurant_is/gen/go/order/v1;orderv1";

import "google/protobuf/empty.proto";
import "google/protobuf/timestamp.proto";

service OrderService {
  // Основные операции с заказами
  rpc CreateOrder(CreateOrderRequest) returns (OrderResponse);
  rpc GetOrder(GetOrderRequest) returns (OrderResponse);
  rpc ListOrders(ListOrdersRequest) returns (ListOrdersResponse);
  rpc UpdateOrderStatus(UpdateOrderStatusRequest) returns (google.protobuf.Empty);

  // Управление элементами заказа
  rpc AddOrderItem(AddOrderItemRequest) returns (OrderItem);
  rpc UpdateOrderItem(UpdateOrderItemRequest) returns (OrderItem);

  // Платежи
  rpc InitiatePayment(PaymentRequest) returns (PaymentResponse);
  rpc ProcessPaymentCallback(PaymentCallbackRequest) returns (google.protobuf.Empty);

  // История и отчетность
  rpc GetOrderHistory(GetOrderRequest) returns (OrderHistoryResponse);
}

// Типы заказов
enum OrderType {
  ORDER_TYPE_UNSPECIFIED = 0;
  ORDER_TYPE_DINE_IN = 1;     // В заведении
  ORDER_TYPE_DELIVERY = 2;    // Доставка
  ORDER_TYPE_TAKEAWAY = 3;    // Самовывоз
}

// Статусы заказа
enum OrderStatus {
  ORDER_STATUS_UNSPECIFIED = 0;
  ORDER_STATUS_CREATED = 1;    // Создан
  ORDER_STATUS_CONFIRMED = 2;  // Подтвержден
  ORDER_STATUS_COOKING = 3;    // Готовится
  ORDER_STATUS_READY = 4;      // Готов к выдаче
  ORDER_STATUS_DELIVERED = 5;  // Доставлен
  ORDER_STATUS_CANCELLED = 6;  // Отменен
}

// Статусы приготовления блюд
enum CookingStatus {
  COOKING_STATUS_UNSPECIFIED = 0;
  COOKING_STATUS_PENDING = 1;   // Ожидает
  COOKING_STATUS_PREPARING = 2; // Готовится
  COOKING_STATUS_READY = 3;     // Готово
  COOKING_STATUS_SERVED = 4;    // Подано
}

message Order {
  string id = 1;                    // UUID
  int64 user_id = 2;                // ID клиента
  int64 staff_id = 3;               // ID сотрудника
  OrderType order_type = 4;
  OrderStatus status = 5;
  string table_number = 6;          // Для DINE_IN
  bytes delivery_address = 7;       // JSONB
  double total_amount = 8;
  google.protobuf.Timestamp created_at = 9;
  google.protobuf.Timestamp updated_at = 10;
  repeated OrderItem items = 11;
}

message OrderItem {
  string item_id = 1;
  int64 dish_id = 2;
  int32 quantity = 3;
  double price = 4;
  string special_requests = 5;
  CookingStatus cooking_status = 6;
  google.protobuf.Timestamp ready_at = 7;
}

message OrderStatusUpdate {
  OrderStatus status = 1;
  string reason = 2;               // Причина изменения
  int64 changed_by = 3;            // ID сотрудника/системы
  google.protobuf.Timestamp changed_at = 4;
}

// Запросы/Ответы

message CreateOrderRequest {
  int64 user_id = 1;
  OrderType order_type = 2;
  oneof location {
    string table_number = 3;       // Для DINE_IN
    bytes delivery_address = 4;    // JSON
  }
  repeated OrderItemCreation items = 5;
}

message OrderItemCreation {
  int64 dish_id = 1;
  int32 quantity = 2;
  string special_requests = 3;
}

message GetOrderRequest {
  string order_id = 1;
}

message ListOrdersRequest {
  optional int64 user_id = 1;      // Фильтр по пользователю
  optional OrderStatus status = 2; // Фильтр по статусу
  int32 page = 3;
  int32 page_size = 4;
}

message ListOrdersResponse {
  repeated Order orders = 1;
  int32 total_count = 2;
}

message UpdateOrderStatusRequest {
  string order_id = 1;
  OrderStatus status = 2;
  string reason = 3;
  int64 changed_by = 4;
}

message AddOrderItemRequest {
  string order_id = 1;
  OrderItemCreation item = 2;
}

message UpdateOrderItemRequest {
  string order_id = 1;
  string item_id = 2;
  optional int32 quantity = 3;
  optional string special_requests = 4;
  optional CookingStatus cooking_status = 5;
}

// Платежи
message PaymentRequest {
  string order_id = 1;
  double amount = 2;
  string payment_method = 3;     // card/cash
  string callback_url = 4;       // Для вебхуков
}

message PaymentResponse {
  string payment_id = 1;
  string payment_url = 2;        // URL для оплаты
  string status = 3;
}

message PaymentCallbackRequest {
  string payment_id = 1;
  string status = 2;             // success/failed
  string transaction_id = 3;
  bytes metadata = 4;            // Дополнительные данные
}

// История
message OrderHistoryResponse {
  repeated OrderStatusUpdate history = 1;
}
```

**PaymentService:**

```proto
syntax = "proto3";

package payment.v1;
option go_package = "github.com/netscrawler/Restaurant_is/gen/go/payment/v1;paymentv1";

import "google/protobuf/empty.proto";
import "google/protobuf/timestamp.proto";

service PaymentService {
  // Основные операции
  rpc CreatePayment(PaymentRequest) returns (PaymentResponse);
  rpc GetPaymentStatus(PaymentStatusRequest) returns (PaymentStatusResponse);
  rpc ProcessRefund(RefundRequest) returns (RefundResponse);

  // Управление поведением заглушки (для тестирования)
  rpc SetMockBehavior(MockConfig) returns (google.protobuf.Empty);
  rpc ResetMockBehavior(google.protobuf.Empty) returns (google.protobuf.Empty);
}

// Статусы платежа
enum PaymentStatus {
  PAYMENT_STATUS_UNSPECIFIED = 0;
  PAYMENT_STATUS_PENDING = 1;     // Ожидает обработки
  PAYMENT_STATUS_SUCCEEDED = 2;   // Успешно завершен
  PAYMENT_STATUS_FAILED = 3;      // Ошибка оплаты
  PAYMENT_STATUS_REFUNDED = 4;    // Возврат средств
}

// Типы платежных методов
enum PaymentMethodType {
  PAYMENT_METHOD_UNSPECIFIED = 0;
  PAYMENT_METHOD_CARD = 1;        // Банковская карта
  PAYMENT_METHOD_CASH = 2;        // Наличные
  PAYMENT_METHOD_APPLE_PAY = 3;
  PAYMENT_METHOD_GOOGLE_PAY = 4;
}

message PaymentRequest {
  string order_id = 1;            // Идентификатор заказа
  double amount = 2;              // Сумма платежа
  PaymentMethodType method = 3;   // Способ оплаты
  map<string, string> metadata = 4; // Дополнительные данные
}

message PaymentResponse {
  string payment_id = 1;          // Уникальный идентификатор платежа
  string status = 2;              // Текущий статус
  string details = 3;             // Детализация (для ошибок)
  google.protobuf.Timestamp created_at = 4;
}

message PaymentStatusRequest {
  string payment_id = 1;
}

message PaymentStatusResponse {
  PaymentStatus status = 1;
  string error_message = 2;       // Причина ошибки, если есть
  google.protobuf.Timestamp updated_at = 3;
}

message RefundRequest {
  string payment_id = 1;
  double amount = 2;              // Сумма для возврата
  string reason = 3;              // Причина возврата
}

message RefundResponse {
  string refund_id = 1;
  PaymentStatus status = 2;
  google.protobuf.Timestamp processed_at = 3;
}

// Конфигурация поведения заглушки
message MockConfig {
  // Процент успешных платежей (0-100)
  int32 success_rate = 1;

  // Фиксированный ответ (если задан, success_rate игнорируется)
  optional PaymentStatus fixed_status = 2;

  // Задержка ответа в миллисекундах
  int32 response_delay_ms = 3;

  // Регулярное выражение для ошибок по метаданным
  string error_pattern = 4;       // Пример: ".*TEST_FAIL.*"

  // Автоматическое время ответа
  bool auto_approve_refunds = 5;  // Авто-подтверждение возвратов
}
```

**GATE:**

```yaml
openapi: 3.0.3
info:
  title: Restaurant Management API
  version: 1.0.0
  description: API Gateway for restaurant management system
servers:
  - url: https://api.your-restaurant.com/api/v1
tags:
  - name: Authentication
    description: User authentication and authorization
  - name: Menu
    description: Menu management
  - name: Orders
    description: Order processing
  - name: Payments
    description: Payment operations

paths:
  /auth/login:
    post:
      tags: [Authentication]
      summary: User login
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: "#/components/schemas/LoginRequest"
      responses:
        200:
          description: Successfully authenticated
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/AuthResponse"
        401:
          $ref: "#/components/responses/Unauthorized"

  /menu/active:
    get:
      tags: [Menu]
      summary: Get active menu
      responses:
        200:
          description: Active menu data
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/MenuResponse"
        404:
          $ref: "#/components/responses/NotFound"

  /orders:
    post:
      tags: [Orders]
      summary: Create new order
      security:
        - BearerAuth: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: "#/components/schemas/OrderRequest"
      responses:
        201:
          description: Order created
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/OrderResponse"
        400:
          $ref: "#/components/responses/BadRequest"

  /orders/{order_id}:
    get:
      tags: [Orders]
      summary: Get order details
      parameters:
        - name: order_id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        200:
          description: Order details
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/Order"
        404:
          $ref: "#/components/responses/NotFound"

components:
  securitySchemes:
    BearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT

  schemas:
    LoginRequest:
      type: object
      properties:
        email:
          type: string
          format: email
        password:
          type: string
          format: password
      required: [email, password]

    AuthResponse:
      type: object
      properties:
        access_token:
          type: string
        expires_in:
          type: integer
        user:
          $ref: "#/components/schemas/User"

    User:
      type: object
      properties:
        id:
          type: integer
        roles:
          type: array
          items:
            type: string
        email:
          type: string

    MenuResponse:
      type: object
      properties:
        categories:
          type: array
          items:
            $ref: "#/components/schemas/Category"

    Category:
      type: object
      properties:
        id:
          type: integer
        name:
          type: string
        dishes:
          type: array
          items:
            $ref: "#/components/schemas/Dish"

    Dish:
      type: object
      properties:
        id:
          type: integer
        name:
          type: string
        price:
          type: number
          format: float
        image_url:
          type: string

    OrderRequest:
      type: object
      properties:
        type:
          type: string
          enum: [dine-in, delivery, takeaway]
        delivery_address:
          type: object
        items:
          type: array
          items:
            $ref: "#/components/schemas/OrderItemRequest"
      required: [type, items]

    OrderItemRequest:
      type: object
      properties:
        dish_id:
          type: integer
        quantity:
          type: integer

    OrderResponse:
      type: object
      properties:
        id:
          type: string
          format: uuid
        status:
          type: string
        total_amount:
          type: number
          format: float

    Order:
      type: object
      properties:
        id:
          type: string
          format: uuid
        status:
          type: string
        items:
          type: array
          items:
            $ref: "#/components/schemas/OrderItem"
        total_amount:
          type: number
          format: float

    OrderItem:
      type: object
      properties:
        dish_id:
          type: integer
        quantity:
          type: integer
        price:
          type: number
          format: float

  responses:
    Unauthorized:
      description: Unauthorized
      content:
        application/json:
          schema:
            type: object
            properties:
              error:
                type: string
                example: "Invalid credentials"

    NotFound:
      description: Resource not found
      content:
        application/json:
          schema:
            type: object
            properties:
              error:
                type: string
                example: "Resource not found"

    BadRequest:
      description: Invalid request
      content:
        application/json:
          schema:
            type: object
            properties:
              error:
                type: string
                example: "Invalid request parameters"
```

=== Структура базы данных

Каждый сервис имеет свою базу данных (PostgreSQL):

**AuthServiceDB:**

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20) UNIQUE,
    password_hash TEXT NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    last_login TIMESTAMPTZ
);
```

```sql
CREATE TABLE staff (
    id SERIAL PRIMARY KEY,
    work_email VARCHAR(255) UNIQUE NOT NULL, -- Отдельный email для работы
    work_phone VARCHAR(20) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL, -- Отдельный хеш пароля
    full_name VARCHAR(100) NOT NULL,
    hire_date DATE NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL, -- 'customer', 'admin', 'chef'
    description TEXT
);
```

PostgreSQL + S3
В PostgreSQL рядом с информацией о меню/блюде хранится идентификатор на s3 по которому можно достать это фото и отдать на фронт
**MenuServiceDB:**

```sql
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,        -- Название категории: "Супы", "Десерты"
    description TEXT,                        -- Описание категории (опционально)
    display_order INT DEFAULT 0,             -- Порядок отображения в меню (0 = первый)
    is_active BOOLEAN DEFAULT TRUE           -- Активна ли категория (скрыть/показать)
);

CREATE INDEX idx_categories_active ON categories(is_active);
CREATE TABLE dishes (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,              -- Название блюда: "Стейк Рибай"
    description TEXT,                        -- Описание состава и особенностей
    price NUMERIC(10, 2) NOT NULL
        CHECK (price > 0),                   -- Цена (например: 1500.50)
    category_id INT NOT NULL
        REFERENCES categories(id)
        ON DELETE RESTRICT,                  -- Запрет удаления категории с блюдами
    cooking_time_min INT
        CHECK (cooking_time_min > 0),        -- Время приготовления в минутах (опционально)
    image_url VARCHAR(500),                  -- Ссылка на фото в S3: "https://bucket.s3.amazonaws.com/dishes/123.jpg"
    is_available BOOLEAN DEFAULT TRUE,       -- Доступно ли для заказа
    calories INT,                            -- Ккал (опционально)
);

CREATE INDEX idx_dishes_category ON dishes(category_id);
CREATE INDEX idx_dishes_availability ON dishes(is_available);

CREATE TABLE promotions (
    id SERIAL PRIMARY KEY,
    dish_id INT NOT NULL REFERENCES dishes(id),
    discount_percent INT CHECK (discount_percent BETWEEN 1 AND 9900), -- Скидка 10%
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    CONSTRAINT valid_dates CHECK (end_date > start_date)
);
CREATE TABLE menus (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,        -- "Основное меню", "Летнее спецпредложение"
    description TEXT,
    valid_from DATE NOT NULL,                 -- Дата начала действия меню
    valid_to DATE NOT NULL,                   -- Дата окончания
    is_active BOOLEAN DEFAULT TRUE,           -- Активно ли для отображения
    created_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT valid_dates CHECK (valid_to >= valid_from)
);
-- Связь меню с категориями
CREATE TABLE menu_categories (
    menu_id INT NOT NULL REFERENCES menus(id) ON DELETE CASCADE,
    category_id INT NOT NULL REFERENCES categories(id),
    display_order INT DEFAULT 0,              -- Порядок категорий внутри меню
    PRIMARY KEY (menu_id, category_id)
);

-- Связь меню с отдельными блюдами (если нужно включать вне категорий)
CREATE TABLE menu_dishes (
    menu_id INT NOT NULL REFERENCES menus(id) ON DELETE CASCADE,
    dish_id INT NOT NULL REFERENCES dishes(id),
    PRIMARY KEY (menu_id, dish_id)
);
```

Эти сервисы будут на Clichouse
**OrderServiceDB:**

```sql
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id INT NOT NULL,                    -- ID клиента (из UserService)
    staff_id INT,                            -- ID сотрудника (официанта, из StaffService)
    type VARCHAR(20) NOT NULL
        CHECK (type IN ('dine-in', 'delivery', 'takeaway')),
    status VARCHAR(50) NOT NULL
        CHECK (status IN ('created', 'confirmed', 'cooking', 'ready', 'delivered', 'canceled')),
    delivery_address JSONB,                  -- Для типа 'delivery' (город, улица, квартира)
    total_amount NUMERIC(10, 2) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    finished_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_user ON orders(user_id);

CREATE TABLE order_items (
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    dish_id INT NOT NULL,                    -- ID блюда (из MenuService)
    quantity INT NOT NULL CHECK (quantity > 0),
    price NUMERIC(10, 2) NOT NULL,           -- Цена на момент заказа (фиксируется)
    special_requests TEXT,                   -- "Без лука, добавить соус"
    cooking_status VARCHAR(50)               -- Статус из KitchenService: 'pending', 'cooking', 'ready'
);

CREATE INDEX idx_order_items_dish ON order_items(dish_id);

CREATE TABLE order_status_history (
    order_id UUID NOT NULL REFERENCES orders(id),
    status VARCHAR(50) NOT NULL,
    changed_by INT,                          -- ID сотрудника или системы (NULL = автоматически)
    changed_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_status_history_order ON order_status_history(order_id);

CREATE TYPE payment_method AS ENUM ('cash', 'card', 'online');

CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL UNIQUE REFERENCES orders(id),
    amount NUMERIC(10, 2) NOT NULL CHECK (amount > 0),
    method payment_method NOT NULL,
    status VARCHAR(20) NOT NULL
        CHECK (status IN ('pending', 'completed', 'failed', 'refunded')),
    transaction_id VARCHAR(255),             -- ID транзакции в payment_service
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_payments_order ON payments(order_id);
```

**PaymentServiceDB:**

```sql
CREATE TYPE payment_status AS ENUM (
    'pending',      -- Платеж инициирован
    'succeeded',    -- Успешно завершен
    'failed',       -- Ошибка оплаты
    'refunded'      -- Средства возвращены
);

CREATE TYPE payment_method_type AS ENUM (
    'card',         -- Банковская карта
    'cash',         -- Наличные (для самовывоза)
);

CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL,                   -- Ссылка на заказ в OrderService
    amount NUMERIC(10, 2) NOT NULL
        CHECK (amount > 0),
    currency VARCHAR(3) NOT NULL
        DEFAULT 'RUB',
    status payment_status NOT NULL
        DEFAULT 'pending',
    error_message TEXT,                       -- Причина ошибки: "Insufficient funds"
    method payment_method_type NOT NULL,
    external_id VARCHAR(255) UNIQUE,          -- ID транзакции в платежном шлюзе (Stripe/YooKassa)
    metadata JSONB,                           -- Доп. данные: {"ip": "192.168.0.1", "user_agent": "..."}
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_payments_order ON payments(order_id);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_external_id ON payments(external_id);


CREATE TABLE payment_methods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id INT NOT NULL,                     -- ID клиента из UserService
    method_type payment_method_type NOT NULL,
    token VARCHAR(500) NOT NULL,              -- Зашифрованные данные карты/кошелька
    is_default BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_payment_methods_user ON payment_methods(user_id);
```

== Implementation

### **Этапы разработки (2 месяца)**

#### **Неделя 1-2: Аутентификация и управление пользователями**

- Разработать **AuthService** (JWT, регистрация, вход).
- Реализовать **UserService** (CRUD пользователей).
- Подключить API Gateway для проверки токенов.
- Начать разработку UI (страницы входа/регистрации).

#### **Неделя 3-4: Управление меню**

- Разработать **MenuService** (CRUD блюд, категории).
- Подключить к фронтенду (страницы меню, редактирование).

#### **Неделя 5-6: Заказы и обработка событий**

- Разработать **OrderService** (создание/обновление заказов).
- Интеграция с Kafka для событий.
- Добавить интерфейсы заказов на фронтенде (клиент, официант, повар).

#### **Неделя 7: Платежи (заглушка) и финальная сборка**

- Разработать **PaymentService** (заглушка).
- Подключить к **OrderService** (изменение статуса заказа после оплаты).
- Завершить разработку UI (страница оплаты, статусы заказов).

#### **Неделя 8: Тестирование, отладка, развертывание**

- Интеграционное тестирование API и UI.
- Настроить мониторинг (Prometheus, Grafana, Loki).
- Настройка деплоя Terraform, Ansible конфигурации, настройка всех сервисов и бд.

== Gathering Results

### **Критерии успешности**

#### **1. Функциональность**

- Все основные сервисы работают корректно.
- Возможность создавать заказы, редактировать меню и управлять пользователями.
- Платёжная заглушка успешно меняет статус заказа.

#### **2. Производительность**

- Система выдерживает **500-1000 одновременных пользователей**.
- Среднее время отклика API **≤ 200 мс**.

#### **3. UI/UX**

- Интерфейс удобен для клиентов, официантов, поваров и администраторов.
- Минимальное количество кликов для оформления заказа.
- Поддержка мобильных устройств.

#### **4. Масштабируемость**

- Возможность интеграции с реальной платёжной системой.
- Готовность к добавлению доставки и других фич.

---
