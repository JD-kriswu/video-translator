#!/usr/bin/env node
/**
 * 抖音视频下载器 - Playwright 版 v2
 * 从 RENDER_DATA 提取真实视频地址
 */

const { chromium } = require('playwright');
const fs = require('fs');
const path = require('path');
const https = require('https');
const http = require('http');

const url = process.argv[2];
const outputPath = process.argv[3] || 'temp/douyin_video.mp4';

if (!url) {
  console.error('用法: node douyin-download.js <抖音链接> [输出路径]');
  process.exit(1);
}

async function downloadFile(fileUrl, dest) {
  return new Promise((resolve, reject) => {
    const file = fs.createWriteStream(dest);
    const protocol = fileUrl.startsWith('https') ? https : http;
    
    const options = {
      headers: {
        'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36',
        'Referer': 'https://www.douyin.com/'
      }
    };
    
    protocol.get(fileUrl, options, (response) => {
      if (response.statusCode === 302 || response.statusCode === 301) {
        file.close();
        fs.unlinkSync(dest);
        downloadFile(response.headers.location, dest).then(resolve).catch(reject);
        return;
      }
      
      if (response.statusCode !== 200) {
        reject(new Error(`下载失败: HTTP ${response.statusCode}`));
        return;
      }
      
      response.pipe(file);
      file.on('finish', () => {
        file.close();
        resolve(dest);
      });
    }).on('error', (err) => {
      fs.unlink(dest, () => {});
      reject(err);
    });
  });
}

async function main() {
  console.error(`正在解析: ${url}`);
  
  const browser = await chromium.launch({
    headless: true,
    args: ['--no-sandbox', '--disable-setuid-sandbox', '--disable-dev-shm-usage']
  });
  
  try {
    const context = await browser.newContext({
      userAgent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36',
      viewport: { width: 1920, height: 1080 },
      locale: 'zh-CN'
    });
    
    await context.addInitScript(() => {
      Object.defineProperty(navigator, 'webdriver', { get: () => false });
    });
    
    const page = await context.newPage();
    
    // 访问页面
    await page.goto(url, { waitUntil: 'domcontentloaded', timeout: 60000 });
    await page.waitForTimeout(5000);
    
    // 从 RENDER_DATA 提取视频信息
    const videoUrl = await page.evaluate(() => {
      const renderData = document.getElementById('RENDER_DATA');
      if (!renderData) return null;
      
      try {
        const data = JSON.parse(decodeURIComponent(renderData.textContent));
        
        // 递归查找视频地址
        const findVideoUrl = (obj, depth = 0) => {
          if (!obj || typeof obj !== 'object' || depth > 15) return null;
          
          // 查找 video 对象下的 playAddr
          if (obj.video && obj.video.playAddr) {
            const playAddr = obj.video.playAddr;
            if (Array.isArray(playAddr) && playAddr.length > 0 && playAddr[0].src) {
              return playAddr[0].src;
            }
          }
          
          // 查找 play_addr
          if (obj.play_addr) {
            if (obj.play_addr.url_list && obj.play_addr.url_list.length > 0) {
              return obj.play_addr.url_list[0];
            }
            if (Array.isArray(obj.play_addr) && obj.play_addr.length > 0) {
              return obj.play_addr[0].src || obj.play_addr[0];
            }
          }
          
          // 查找 playApi
          if (obj.playApi) {
            return obj.playApi;
          }
          
          // 递归搜索
          for (const key of Object.keys(obj)) {
            // 跳过特效相关的字段
            if (key === 'effect' || key === 'sticker' || key === 'effectcdn') continue;
            
            const result = findVideoUrl(obj[key], depth + 1);
            if (result && !result.includes('effectcdn') && !result.includes('ies.fe.effect')) {
              return result;
            }
          }
          return null;
        };
        
        return findVideoUrl(data);
      } catch (e) {
        console.error('解析 RENDER_DATA 失败:', e);
        return null;
      }
    });
    
    // 如果 RENDER_DATA 没找到，尝试从网络请求捕获
    let finalUrl = videoUrl;
    
    if (!finalUrl) {
      console.error('RENDER_DATA 未找到视频，尝试从页面元素获取...');
      
      // 尝试点击播放
      await page.click('video').catch(() => {});
      await page.waitForTimeout(3000);
      
      finalUrl = await page.evaluate(() => {
        // 从 xg-video-container 获取
        const container = document.querySelector('xg-video-container video');
        if (container && container.src && !container.src.startsWith('blob:')) {
          return container.src;
        }
        
        // 从所有 video 标签获取
        const videos = document.querySelectorAll('video');
        for (const v of videos) {
          if (v.src && !v.src.startsWith('blob:') && !v.src.includes('effectcdn')) {
            return v.src;
          }
          const source = v.querySelector('source');
          if (source && source.src && !source.src.includes('effectcdn')) {
            return source.src;
          }
        }
        
        return null;
      });
    }
    
    if (!finalUrl) {
      throw new Error('无法获取视频地址');
    }
    
    // 处理地址
    finalUrl = finalUrl.replace(/\\u002F/g, '/');
    // 去水印
    finalUrl = finalUrl.replace('playwm', 'play');
    
    // 验证不是特效素材
    if (finalUrl.includes('effectcdn') || finalUrl.includes('ies.fe.effect')) {
      throw new Error('获取到的是特效素材，不是视频内容');
    }
    
    console.error(`视频地址: ${finalUrl.substring(0, 80)}...`);
    
    // 确保输出目录存在
    const dir = path.dirname(outputPath);
    if (!fs.existsSync(dir)) {
      fs.mkdirSync(dir, { recursive: true });
    }
    
    // 下载视频
    console.error('正在下载视频...');
    await downloadFile(finalUrl, outputPath);
    
    // 检查文件大小
    const stats = fs.statSync(outputPath);
    if (stats.size < 500 * 1024) { // 小于 500KB 可能不是完整视频
      console.error(`警告: 文件较小 (${Math.round(stats.size/1024)}KB)，可能不是完整视频`);
    }
    
    console.log(outputPath);
    
  } finally {
    await browser.close();
  }
}

main().catch(err => {
  console.error(`错误: ${err.message}`);
  process.exit(1);
});
