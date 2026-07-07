import { createCliRenderer } from "@opentui/core";
import { createRoot } from "@opentui/react";
import { App } from "./app";
import { connect } from "./client";

const noWatch = process.argv.includes("--no-watch");
const client = connect({ noWatch });

const renderer = await createCliRenderer({
  exitOnCtrlC: false, // we quit via our own handler so the terminal is always restored
  consoleMode: "disabled",
});

let quitting = false;
function quit(code = 0): void {
  if (quitting) return;
  quitting = true;
  client.close();
  renderer.destroy();
  process.exit(code);
}

process.on("SIGINT", () => quit(130));
process.on("SIGTERM", () => quit(143));

createRoot(renderer).render(<App client={client} quit={() => quit()} />);
