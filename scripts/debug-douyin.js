#!/usr/bin/env node
// 调试脚本 - 打印 RENDER_DATA 内容

const { chromium } = require('playwright');

const url = process.argv[2] || 'https://www.douyin.com/video/7597354815329011822';

async function main() {
  const browser = await chromium.launch({ headless: true, args: ['--no-sandbox'] });
  const context = await browser.newContext({
    userAgent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36',
    viewport: { width: 1920, height: 1080 }
  });
  
  const page = await context.newPage();
  await page.goto(url, { waitUntil: 'domcontentloaded', timeout: 60000 });
  await page.waitForTimeout(5000);
  
  const data = await page.evaluate(() => {
    const el = document.getElementById('RENDER_DATA');
    if (!el) return 'RENDER_DATA not found';
    return decodeURIComponent(el.textContent);
  });
  
  console.log(data);
  await browser.close();
}

main().catch(console.error);
