import "@typespec/http";
import "@typespec/versioning";

using Http;
using Versioning;

@service(#{ title: "controlriver.com" })
@server("http://localhost:8080/api/v1")
@versioned(Versions)
namespace ControlRiver;

enum Versions {
  v1: "1.0",
}

@added(Versions.v1)
model ErrorResponse {
  error: string;
  message: string;
}

@added(Versions.v1)
model HealthResponse {
  message: string;
  status: string;
  timestamp: string;
}

@added(Versions.v1)
model CreateAccountRequest {
  name: string;
  description?: string;
  currency_code: string;
}

@added(Versions.v1)
model CreateAccountResponse {
  account_id: string;
}

@route("/health")
@tag("Health")
namespace Health {
  @added(Versions.v1)
  @get
  op getHealth(): {
    @statusCode statusCode: 200;
    @body health: HealthResponse;
  };
}

@route("/users")
@tag("Users")
namespace Users {
  @added(Versions.v1)
  @post
  op createUser(
    @body request: {
      name: string;
      email: string;
    },
  ): {
    @statusCode statusCode: 201;
  } | {
    @statusCode statusCode: 400;
    @body error: ErrorResponse;
  } | {
    @statusCode statusCode: 500;
    @body error: ErrorResponse;
  };
}

@route("/accounts")
@tag("Accounts")
namespace Accounts {
  @added(Versions.v1)
  @post
  op createAccount(
    @header("X-User-ID") userId: string,
    @body request: CreateAccountRequest,
  ): {
    @statusCode statusCode: 201;
    @body response: CreateAccountResponse;
  } | {
    @statusCode statusCode: 400;
    @body error: string[];
  } | {
    @statusCode statusCode: 500;
    @body error: string[];
  };
}

@route("/categories")
@tag("Categories")
namespace Categories {
  @added(Versions.v1)
  @post
  op createCategory(
    @header("X-Account-ID") accountId: string,
    @body request: {
      name: string;
      type: "expense" | "income";

      @example("#FF5733")
      color: string;

      @example("AC-09")
      parent_id?: string;
    },
  ): {
    @statusCode statusCode: 201;
    @body response: "Category created successfully";
  } | {
    @statusCode statusCode: 400;
    @body error: ErrorResponse;
  } | {
    @statusCode statusCode: 500;
    @body error: ErrorResponse;
  };
}

@route("/transactions")
@tag("Transactions")
namespace Transactions {
  @added(Versions.v1)
  @post
  op createTransaction(
    @header("X-Account-ID") accountId: string,
    @body request: {
      @example("AC-09")
      category_code: string;

      amount: numeric;
      description?: string;

      @example("2023-10-01")
      transaction_date: string;
    },
  ): {
    @statusCode statusCode: 201;
    @body response: "Transaction created successfully";
  } | {
    @statusCode statusCode: 400;
    @body error: ErrorResponse;
  } | {
    @statusCode statusCode: 500;
    @body error: ErrorResponse;
  };
}

@added(Versions.v1)
model TransactionReport {
  amount: numeric;
  description?: string;
  date: string;
}

@added(Versions.v1)
model CategoryReport {
  name: string;
  total_spent: numeric;
  transactions?: TransactionReport[];
  sub_categories?: {
    name: string;
    total_spent: numeric;
    transactions?: TransactionReport[];
  }[];
}

@tag("Reports")
@route("/reports/by-category")
namespace Reports {
  @added(Versions.v1)
  @get
  op getByCategoryReport(
    @header("X-Account-ID") accountId: string,
    @query("type") type: "expense" | "income",
    @query("from") from: string,
    @query("to") to: string,
    @query("with-sub-categories") withSubCategories: boolean = false,
    @query("with-transactions") withTransactions: boolean = false,
  ): {
    @statusCode statusCode: 200;
    @body response: {
      from: string;
      to: string;
      type: "expense" | "income";
      categories: CategoryReport[];
    };
  } | {
    @statusCode statusCode: 400;
    @body error: ErrorResponse;
  } | {
    @statusCode statusCode: 500;
    @body error: ErrorResponse;
  };
}
