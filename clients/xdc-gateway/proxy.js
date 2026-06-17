const http = require('http');
const { createProxyServer } = require('http-proxy');

const UPSTREAM = process.env.UPSTREAM_URL || 'http://127.0.0.1:8546';
const ENGINE_UPSTREAM = process.env.ENGINE_URL || 'http://127.0.0.1:8651';
const PORT = process.env.PORT || 8545;
const ENGINE_PORT = process.env.ENGINE_PORT || 8551;

const SECRET = Buffer.from('secretsecretsecretsecretsecretsecretsecret', 'utf8').slice(0, 32);

function createJwt(secret) {
  const header = Buffer.from(JSON.stringify({ alg: 'HS256', typ: 'JWT' })).toString('base64url');
  const now = Math.floor(Date.now() / 1000);
  const payload = Buffer.from(JSON.stringify({ iat: now })).toString('base64url');
  const data = `${header}.${payload}`;
  const sig = require('crypto').createHmac('sha256', secret).update(data).digest('base64url');
  return `${data}.${sig}`;
}

const proxy = createProxyServer({
  target: UPSTREAM,
  changeOrigin: true,
  ws: true,
});

const engineProxy = createProxyServer({
  target: ENGINE_UPSTREAM,
  changeOrigin: true,
});

function handleError(err, req, res) {
  console.error('Proxy error:', err.message);
  if (res && !res.headersSent) {
    res.writeHead(502, { 'Content-Type': 'application/json' });
    res.end(JSON.stringify({ jsonrpc: '2.0', error: { code: -32603, message: 'Gateway proxy error: ' + err.message }, id: null }));
  }
}

proxy.on('error', handleError);
engineProxy.on('error', handleError);

const server = http.createServer((req, res) => {
  console.log(`[gateway] ${req.method} ${req.url}`);
  proxy.web(req, res);
});

server.on('upgrade', (req, socket, head) => {
  proxy.ws(req, socket, head);
});

const engineServer = http.createServer((req, res) => {
  console.log(`[gateway-engine] ${req.method} ${req.url}`);
  req.headers['authorization'] = 'Bearer ' + createJwt(SECRET);
  engineProxy.web(req, res);
});

engineServer.on('upgrade', (req, socket, head) => {
  req.headers['authorization'] = 'Bearer ' + createJwt(SECRET);
  engineProxy.ws(req, socket, head);
});

server.listen(PORT, '0.0.0.0', () => {
  console.log(`XDC Gateway proxy listening on 0.0.0.0:${PORT} -> ${UPSTREAM}`);
});

engineServer.listen(ENGINE_PORT, '0.0.0.0', () => {
  console.log(`XDC Gateway engine proxy listening on 0.0.0.0:${ENGINE_PORT} -> ${ENGINE_UPSTREAM}`);
});
