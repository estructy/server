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
