import Chrome from 'puppeteer';
import { Context } from 'uvu';

// Launch the browser
// Add `browser` and `page` to context
export async function setup(context: Context) {
	context.browser = await Chrome.launch();
	context.page = await context.browser.newPage();
}

// Close everything on suite completion
export async function reset(context: Context) {
	await context.page?.close();
	await context.browser?.close();
}

// Navigate to homepage
export async function homepage(context: Context) {
	await context.page.goto('http://localhost:8082');
}