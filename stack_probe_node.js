// file: stack_probe_node.js
let depth = 0;

function recurse(n) {
  const buf = new Array(1024).fill(n);
  depth = n;
  recurse(n + 1);
}

console.log("start node recurse test");

try {
  recurse(1);
} catch (e) {
  const mu = process.memoryUsage();
  console.error(
    `>>> CRASH at depth=${depth} | rss=${Math.round(mu.rss / 1024)}KB | heapUsed=${Math.round(
      mu.heapUsed / 1024
    )}KB | external=${Math.round(mu.external / 1024)}KB`
  );
  console.error(e.message);
  process.exit(1);
}

