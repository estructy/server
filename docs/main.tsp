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

@route("/accounts")
@tag("Accounts")
namespace Accounts {
  @added(Versions.v1)
  @post
	op createAccount(@header("X-User-ID") userId: string, @body request: CreateAccountRequest): {
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
