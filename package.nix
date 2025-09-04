{
  lib,
  buildGoModule,
}:
buildGoModule {
  name = "Redes_RemoteMCPServer";
  src = ./.;
  # vendorHash = "sha256-fTP/PZXcJUuDx3OA2zJSTqGTwcIAJI7qXeWlCit9f+k=";
  vendorHash = "sha256-8Vs08PkqO/voV08HThCr7lwQ7AeEa501qhgBTwWlxB4=";
  meta = {
    description = "FAGD MCP Server for Redes course";
    homepage = "https://github.com/ElrohirGT/Redes_MCPServer";
    license = lib.licenses.mit;
    maintainers = with lib.maintainers; [elrohirgt];
  };
}
