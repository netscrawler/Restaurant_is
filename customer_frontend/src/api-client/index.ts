/* eslint-disable */
/* tslint:disable */
// @ts-nocheck
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

export interface User {
  id?: string;
  email?: string;
  phone?: string;
  name?: string;
  roles?: string[];
  /** @format date-time */
  created_at?: string;
  /** @format date-time */
  updated_at?: string;
}

export interface Staff {
  id?: string;
  work_email?: string;
  position?: string;
  roles?: string[];
  /** @format date-time */
  created_at?: string;
  /** @format date-time */
  updated_at?: string;
}

export interface Dish {
  id?: string;
  name?: string;
  description?: string;
  price?: number;
  category_id?: number;
  category_name?: string;
  image_url?: string;
  available?: boolean;
  /** @format date-time */
  created_at?: string;
  /** @format date-time */
  updated_at?: string;
}

export interface Order {
  id?: string;
  user_id?: string;
  items?: {
    dish_id?: string;
    dish_name?: string;
    quantity?: number;
    price?: number;
  }[];
  total_amount?: number;
  status?:
    | "ORDER_STATUS_CREATED"
    | "ORDER_STATUS_CONFIRMED"
    | "ORDER_STATUS_COOKING"
    | "ORDER_STATUS_READY"
    | "ORDER_STATUS_DELIVERED"
    | "ORDER_STATUS_CANCELLED";
  delivery_address?: string;
  comment?: string;
  /** @format date-time */
  created_at?: string;
  /** @format date-time */
  updated_at?: string;
}

export interface Error {
  code?: number;
  message?: string;
  details?: string[];
}

export type QueryParamsType = Record<string | number, any>;
export type ResponseFormat = keyof Omit<Body, "body" | "bodyUsed">;

export interface FullRequestParams extends Omit<RequestInit, "body"> {
  /** set parameter to `true` for call `securityWorker` for this request */
  secure?: boolean;
  /** request path */
  path: string;
  /** content type of request body */
  type?: ContentType;
  /** query params */
  query?: QueryParamsType;
  /** format of response (i.e. response.json() -> format: "json") */
  format?: ResponseFormat;
  /** request body */
  body?: unknown;
  /** base url */
  baseUrl?: string;
  /** request cancellation token */
  cancelToken?: CancelToken;
}

export type RequestParams = Omit<
  FullRequestParams,
  "body" | "method" | "query" | "path"
>;

export interface ApiConfig<SecurityDataType = unknown> {
  baseUrl?: string;
  baseApiParams?: Omit<RequestParams, "baseUrl" | "cancelToken" | "signal">;
  securityWorker?: (
    securityData: SecurityDataType | null,
  ) => Promise<RequestParams | void> | RequestParams | void;
  customFetch?: typeof fetch;
}

export interface HttpResponse<D extends unknown, E extends unknown = unknown>
  extends Response {
  data: D;
  error: E;
}

type CancelToken = Symbol | string | number;

export enum ContentType {
  Json = "application/json",
  JsonApi = "application/vnd.api+json",
  FormData = "multipart/form-data",
  UrlEncoded = "application/x-www-form-urlencoded",
  Text = "text/plain",
}

export class HttpClient<SecurityDataType = unknown> {
  public baseUrl: string = "http://localhost:8080/api/v1";
  private securityData: SecurityDataType | null = null;
  private securityWorker?: ApiConfig<SecurityDataType>["securityWorker"];
  private abortControllers = new Map<CancelToken, AbortController>();
  private customFetch = (...fetchParams: Parameters<typeof fetch>) =>
    fetch(...fetchParams);

  private baseApiParams: RequestParams = {
    credentials: "same-origin",
    headers: {},
    redirect: "follow",
    referrerPolicy: "no-referrer",
  };

  constructor(apiConfig: ApiConfig<SecurityDataType> = {}) {
    Object.assign(this, apiConfig);
  }

  public setSecurityData = (data: SecurityDataType | null) => {
    this.securityData = data;
  };

  protected encodeQueryParam(key: string, value: any) {
    const encodedKey = encodeURIComponent(key);
    return `${encodedKey}=${encodeURIComponent(typeof value === "number" ? value : `${value}`)}`;
  }

  protected addQueryParam(query: QueryParamsType, key: string) {
    return this.encodeQueryParam(key, query[key]);
  }

  protected addArrayQueryParam(query: QueryParamsType, key: string) {
    const value = query[key];
    return value.map((v: any) => this.encodeQueryParam(key, v)).join("&");
  }

  protected toQueryString(rawQuery?: QueryParamsType): string {
    const query = rawQuery || {};
    const keys = Object.keys(query).filter(
      (key) => "undefined" !== typeof query[key],
    );
    return keys
      .map((key) =>
        Array.isArray(query[key])
          ? this.addArrayQueryParam(query, key)
          : this.addQueryParam(query, key),
      )
      .join("&");
  }

  protected addQueryParams(rawQuery?: QueryParamsType): string {
    const queryString = this.toQueryString(rawQuery);
    return queryString ? `?${queryString}` : "";
  }

  private contentFormatters: Record<ContentType, (input: any) => any> = {
    [ContentType.Json]: (input: any) =>
      input !== null && (typeof input === "object" || typeof input === "string")
        ? JSON.stringify(input)
        : input,
    [ContentType.JsonApi]: (input: any) =>
      input !== null && (typeof input === "object" || typeof input === "string")
        ? JSON.stringify(input)
        : input,
    [ContentType.Text]: (input: any) =>
      input !== null && typeof input !== "string"
        ? JSON.stringify(input)
        : input,
    [ContentType.FormData]: (input: any) =>
      Object.keys(input || {}).reduce((formData, key) => {
        const property = input[key];
        formData.append(
          key,
          property instanceof Blob
            ? property
            : typeof property === "object" && property !== null
              ? JSON.stringify(property)
              : `${property}`,
        );
        return formData;
      }, new FormData()),
    [ContentType.UrlEncoded]: (input: any) => this.toQueryString(input),
  };

  protected mergeRequestParams(
    params1: RequestParams,
    params2?: RequestParams,
  ): RequestParams {
    return {
      ...this.baseApiParams,
      ...params1,
      ...(params2 || {}),
      headers: {
        ...(this.baseApiParams.headers || {}),
        ...(params1.headers || {}),
        ...((params2 && params2.headers) || {}),
      },
    };
  }

  protected createAbortSignal = (
    cancelToken: CancelToken,
  ): AbortSignal | undefined => {
    if (this.abortControllers.has(cancelToken)) {
      const abortController = this.abortControllers.get(cancelToken);
      if (abortController) {
        return abortController.signal;
      }
      return void 0;
    }

    const abortController = new AbortController();
    this.abortControllers.set(cancelToken, abortController);
    return abortController.signal;
  };

  public abortRequest = (cancelToken: CancelToken) => {
    const abortController = this.abortControllers.get(cancelToken);

    if (abortController) {
      abortController.abort();
      this.abortControllers.delete(cancelToken);
    }
  };

  public request = async <T = any, E = any>({
    body,
    secure,
    path,
    type,
    query,
    format,
    baseUrl,
    cancelToken,
    ...params
  }: FullRequestParams): Promise<HttpResponse<T, E>> => {
    const secureParams =
      ((typeof secure === "boolean" ? secure : this.baseApiParams.secure) &&
        this.securityWorker &&
        (await this.securityWorker(this.securityData))) ||
      {};
    const requestParams = this.mergeRequestParams(params, secureParams);
    const queryString = query && this.toQueryString(query);
    const payloadFormatter = this.contentFormatters[type || ContentType.Json];
    const responseFormat = format || requestParams.format;

    return this.customFetch(
      `${baseUrl || this.baseUrl || ""}${path}${queryString ? `?${queryString}` : ""}`,
      {
        ...requestParams,
        headers: {
          ...(requestParams.headers || {}),
          ...(type && type !== ContentType.FormData
            ? { "Content-Type": type }
            : {}),
        },
        signal:
          (cancelToken
            ? this.createAbortSignal(cancelToken)
            : requestParams.signal) || null,
        body:
          typeof body === "undefined" || body === null
            ? null
            : payloadFormatter(body),
      },
    ).then(async (response) => {
      const r = response.clone() as HttpResponse<T, E>;
      r.data = null as unknown as T;
      r.error = null as unknown as E;

      const data = !responseFormat
        ? r
        : await response[responseFormat]()
            .then((data) => {
              if (r.ok) {
                r.data = data;
              } else {
                r.error = data;
              }
              return r;
            })
            .catch((e) => {
              r.error = e;
              return r;
            });

      if (cancelToken) {
        this.abortControllers.delete(cancelToken);
      }

      if (!response.ok) throw data;
      return data;
    });
  };
}

/**
 * @title Restaurant API Gateway - Combined API
 * @version 1.0.0
 * @license Apache 2.0 (http://www.apache.org/licenses/LICENSE-2.0.html)
 * @baseUrl http://localhost:8080/api/v1
 * @contact API Support <support@swagger.io> (http://www.swagger.io/support)
 *
 * Объединенный API для ресторанной системы. Включает в себя все сервисы: Auth, User, Menu, Order, Notify
 */
export class Api<
  SecurityDataType extends unknown,
> extends HttpClient<SecurityDataType> {
  health = {
    /**
     * @description Проверка состояния сервиса
     *
     * @tags Health
     * @name HealthCheck
     * @summary Health check
     * @request GET:/health
     */
    healthCheck: (params: RequestParams = {}) =>
      this.request<
        {
          /** @example "ok" */
          status?: string;
          /** @example "gate" */
          service?: string;
          /** @example "1.0.0" */
          version?: string;
        },
        any
      >({
        path: `/health`,
        method: "GET",
        format: "json",
        ...params,
      }),
  };
  auth = {
    /**
     * @description Отправляет код подтверждения на телефон клиента
     *
     * @tags Auth
     * @name LoginClientInit
     * @summary Инициализация входа клиента
     * @request POST:/auth/client/login/init
     */
    loginClientInit: (
      data: {
        /** Номер телефона клиента */
        phone: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        {
          status?: string;
          error?: string;
        },
        any
      >({
        path: `/auth/client/login/init`,
        method: "POST",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Подтверждает вход клиента с помощью кода
     *
     * @tags Auth
     * @name LoginClientConfirm
     * @summary Подтверждение входа клиента
     * @request POST:/auth/client/login/confirm
     */
    loginClientConfirm: (
      data: {
        /** Номер телефона клиента */
        phone: string;
        /** Код подтверждения */
        code: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        {
          access_token?: string;
          expires_in?: number;
          refresh_token?: string;
          refresh_token_expires_in?: number;
          user?: User;
        },
        any
      >({
        path: `/auth/client/login/confirm`,
        method: "POST",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Вход сотрудника с помощью email и пароля
     *
     * @tags Auth
     * @name LoginStaff
     * @summary Вход сотрудника
     * @request POST:/auth/staff/login
     */
    loginStaff: (
      data: {
        staff: Staff;
        /** Пароль сотрудника */
        password: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        {
          access_token?: string;
          expires_in?: number;
          refresh_token?: string;
          refresh_token_expires_in?: number;
          user?: User;
        },
        any
      >({
        path: `/auth/staff/login`,
        method: "POST",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Регистрирует нового сотрудника
     *
     * @tags Auth
     * @name RegisterStaff
     * @summary Регистрация сотрудника
     * @request POST:/auth/staff/register
     */
    registerStaff: (
      data: {
        staff: Staff;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        {
          staff?: Staff;
        },
        any
      >({
        path: `/auth/staff/register`,
        method: "POST",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Вход через OAuth Яндекс
     *
     * @tags Auth
     * @name LoginYandex
     * @summary Вход через Яндекс
     * @request POST:/auth/yandex/login
     */
    loginYandex: (
      data: {
        /** Код авторизации от Яндекс */
        code: string;
        /** URI перенаправления */
        redirect_uri: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        {
          access_token?: string;
          expires_in?: number;
          refresh_token?: string;
          refresh_token_expires_in?: number;
          user?: User;
        },
        any
      >({
        path: `/auth/yandex/login`,
        method: "POST",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Обновляет access token с помощью refresh token
     *
     * @tags Auth
     * @name Refresh
     * @summary Обновление токена
     * @request POST:/auth/refresh
     */
    refresh: (
      data: {
        /** Refresh token */
        refresh_token: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        {
          access_token?: string;
          expires_in?: number;
          refresh_token?: string;
          refresh_token_expires_in?: number;
          user?: User;
        },
        any
      >({
        path: `/auth/refresh`,
        method: "POST",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Проверяет валидность токена и возвращает информацию о пользователе
     *
     * @tags Auth
     * @name Validate
     * @summary Валидация токена
     * @request POST:/auth/validate
     */
    validate: (
      data: {
        /** JWT токен для валидации */
        token: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        {
          valid?: boolean;
          user?: User;
        },
        any
      >({
        path: `/auth/validate`,
        method: "POST",
        body: data,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),
  };
  menu = {
    /**
     * @description Возвращает список блюд с возможностью фильтрации
     *
     * @tags Menu
     * @name ListDishes
     * @summary Получить список блюд
     * @request GET:/menu/dishes
     */
    listDishes: (
      query?: {
        /** ID категории для фильтрации */
        category_id?: number;
        /** Только доступные блюда */
        only_available?: boolean;
        /** Номер страницы */
        page?: number;
        /** Размер страницы */
        page_size?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        {
          dishes?: Dish[];
          total?: number;
          page?: number;
          page_size?: number;
        },
        any
      >({
        path: `/menu/dishes`,
        method: "GET",
        query: query,
        format: "json",
        ...params,
      }),

    /**
     * @description Возвращает информацию о конкретном блюде
     *
     * @tags Menu
     * @name GetDish
     * @summary Получить блюдо по ID
     * @request GET:/menu/dishes/{id}
     */
    getDish: (id: string, params: RequestParams = {}) =>
      this.request<Dish, any>({
        path: `/menu/dishes/${id}`,
        method: "GET",
        format: "json",
        ...params,
      }),
  };
  users = {
    /**
     * @description Возвращает профиль текущего пользователя
     *
     * @tags User
     * @name GetProfile
     * @summary Получить профиль пользователя
     * @request GET:/users/profile
     * @secure
     */
    getProfile: (params: RequestParams = {}) =>
      this.request<User, any>({
        path: `/users/profile`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * @description Обновляет профиль текущего пользователя
     *
     * @tags User
     * @name UpdateProfile
     * @summary Обновить профиль пользователя
     * @request PUT:/users/profile
     * @secure
     */
    updateProfile: (
      data: {
        email?: string;
        phone?: string;
        name?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<User, any>({
        path: `/users/profile`,
        method: "PUT",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),
  };
  orders = {
    /**
     * @description Возвращает список заказов пользователя
     *
     * @tags Order
     * @name ListOrders
     * @summary Получить список заказов
     * @request GET:/orders
     * @secure
     */
    listOrders: (
      query?: {
        /** Статус заказа для фильтрации */
        status?:
          | "ORDER_STATUS_CREATED"
          | "ORDER_STATUS_CONFIRMED"
          | "ORDER_STATUS_COOKING"
          | "ORDER_STATUS_READY"
          | "ORDER_STATUS_DELIVERED"
          | "ORDER_STATUS_CANCELLED";
        /** Номер страницы */
        page?: number;
        /** Размер страницы */
        page_size?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        {
          orders?: Order[];
          total?: number;
          page?: number;
          page_size?: number;
        },
        any
      >({
        path: `/orders`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * @description Создает новый заказ
     *
     * @tags Order
     * @name CreateOrder
     * @summary Создать заказ
     * @request POST:/orders
     * @secure
     */
    createOrder: (
      data: {
        items: {
          dish_id?: string;
          quantity?: number;
        }[];
        delivery_address?: string;
        comment?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<Order, any>({
        path: `/orders`,
        method: "POST",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Возвращает информацию о конкретном заказе
     *
     * @tags Order
     * @name GetOrder
     * @summary Получить заказ по ID
     * @request GET:/orders/{id}
     * @secure
     */
    getOrder: (id: string, params: RequestParams = {}) =>
      this.request<Order, any>({
        path: `/orders/${id}`,
        method: "GET",
        secure: true,
        format: "json",
        ...params,
      }),
  };
  admin = {
    /**
     * @description Возвращает список всех пользователей (требует роль admin)
     *
     * @tags Admin
     * @name ListUsers
     * @summary Получить список всех пользователей
     * @request GET:/admin/users
     * @secure
     */
    listUsers: (
      query?: {
        /** Только активные пользователи */
        only_active?: boolean;
        /** Номер страницы */
        page?: number;
        /** Размер страницы */
        page_size?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        {
          users?: User[];
          total?: number;
        },
        any
      >({
        path: `/admin/users`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * @description Создает нового пользователя (требует роль admin)
     *
     * @tags Admin
     * @name CreateUser
     * @summary Создать пользователя
     * @request POST:/admin/users
     * @secure
     */
    createUser: (
      data: {
        email?: string;
        phone?: string;
        name?: string;
        roles?: string[];
      },
      params: RequestParams = {},
    ) =>
      this.request<User, any>({
        path: `/admin/users`,
        method: "POST",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Обновляет информацию о пользователе (требует роль admin)
     *
     * @tags Admin
     * @name UpdateUser
     * @summary Обновить пользователя
     * @request PUT:/admin/users/{id}
     * @secure
     */
    updateUser: (
      id: string,
      data: {
        email?: string;
        phone?: string;
        name?: string;
        roles?: string[];
      },
      params: RequestParams = {},
    ) =>
      this.request<User, any>({
        path: `/admin/users/${id}`,
        method: "PUT",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Удаляет пользователя (требует роль admin)
     *
     * @tags Admin
     * @name DeleteUser
     * @summary Удалить пользователя
     * @request DELETE:/admin/users/{id}
     * @secure
     */
    deleteUser: (id: string, params: RequestParams = {}) =>
      this.request<void, any>({
        path: `/admin/users/${id}`,
        method: "DELETE",
        secure: true,
        ...params,
      }),

    /**
     * @description Возвращает список всех сотрудников (требует роль admin)
     *
     * @tags Admin
     * @name ListStaff
     * @summary Получить список сотрудников
     * @request GET:/admin/staff
     * @secure
     */
    listStaff: (
      query?: {
        /** Только активные сотрудники */
        only_active?: boolean;
        /** Номер страницы */
        page?: number;
        /** Размер страницы */
        page_size?: number;
      },
      params: RequestParams = {},
    ) =>
      this.request<
        {
          staff?: Staff[];
          total?: number;
        },
        any
      >({
        path: `/admin/staff`,
        method: "GET",
        query: query,
        secure: true,
        format: "json",
        ...params,
      }),

    /**
     * @description Обновляет информацию о сотруднике (требует роль admin)
     *
     * @tags Admin
     * @name UpdateStaff
     * @summary Обновить сотрудника
     * @request PUT:/admin/staff/{id}
     * @secure
     */
    updateStaff: (
      id: string,
      data: {
        work_email?: string;
        position?: string;
        password?: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<Staff, any>({
        path: `/admin/staff/${id}`,
        method: "PUT",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Назначает роль сотруднику (требует роль admin)
     *
     * @tags Admin
     * @name AssignRole
     * @summary Назначить роль сотруднику
     * @request POST:/admin/staff/{id}/roles
     * @secure
     */
    assignRole: (
      id: string,
      data: {
        role: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<void, any>({
        path: `/admin/staff/${id}/roles`,
        method: "POST",
        body: data,
        secure: true,
        type: ContentType.Json,
        ...params,
      }),

    /**
     * @description Отзывает роль у сотрудника (требует роль admin)
     *
     * @tags Admin
     * @name RevokeRole
     * @summary Отозвать роль у сотрудника
     * @request DELETE:/admin/staff/{id}/roles
     * @secure
     */
    revokeRole: (
      id: string,
      data: {
        role: string;
      },
      params: RequestParams = {},
    ) =>
      this.request<void, any>({
        path: `/admin/staff/${id}/roles`,
        method: "DELETE",
        body: data,
        secure: true,
        type: ContentType.Json,
        ...params,
      }),

    /**
     * @description Создает новое блюдо (требует роль admin)
     *
     * @tags Admin
     * @name CreateDish
     * @summary Создать блюдо
     * @request POST:/admin/menu/dishes
     * @secure
     */
    createDish: (
      data: {
        name: string;
        description?: string;
        price: number;
        category_id?: number;
        image_url?: string;
        available?: boolean;
      },
      params: RequestParams = {},
    ) =>
      this.request<Dish, any>({
        path: `/admin/menu/dishes`,
        method: "POST",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Обновляет информацию о блюде (требует роль admin)
     *
     * @tags Admin
     * @name UpdateDish
     * @summary Обновить блюдо
     * @request PUT:/admin/menu/dishes/{id}
     * @secure
     */
    updateDish: (
      id: string,
      data: {
        name?: string;
        description?: string;
        price?: number;
        category_id?: number;
        image_url?: string;
        available?: boolean;
      },
      params: RequestParams = {},
    ) =>
      this.request<Dish, any>({
        path: `/admin/menu/dishes/${id}`,
        method: "PUT",
        body: data,
        secure: true,
        type: ContentType.Json,
        format: "json",
        ...params,
      }),

    /**
     * @description Удаляет блюдо (требует роль admin)
     *
     * @tags Admin
     * @name DeleteDish
     * @summary Удалить блюдо
     * @request DELETE:/admin/menu/dishes/{id}
     * @secure
     */
    deleteDish: (id: string, params: RequestParams = {}) =>
      this.request<void, any>({
        path: `/admin/menu/dishes/${id}`,
        method: "DELETE",
        secure: true,
        ...params,
      }),

    /**
     * @description Обновляет статус заказа (требует роль admin)
     *
     * @tags Admin
     * @name UpdateOrderStatus
     * @summary Обновить статус заказа
     * @request PUT:/admin/orders/{id}/status
     * @secure
     */
    updateOrderStatus: (
      id: string,
      data: {
        status:
          | "ORDER_STATUS_CREATED"
          | "ORDER_STATUS_CONFIRMED"
          | "ORDER_STATUS_COOKING"
          | "ORDER_STATUS_READY"
          | "ORDER_STATUS_DELIVERED"
          | "ORDER_STATUS_CANCELLED";
      },
      params: RequestParams = {},
    ) =>
      this.request<void, any>({
        path: `/admin/orders/${id}/status`,
        method: "PUT",
        body: data,
        secure: true,
        type: ContentType.Json,
        ...params,
      }),
  };
}
