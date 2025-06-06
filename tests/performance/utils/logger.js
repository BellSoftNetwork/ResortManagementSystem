/**
 * k6 테스트용 로깅 유틸리티
 */

export class Logger {
  static info(message, data = null) {
    console.log(`[정보] ${message}`, data ? JSON.stringify(data) : '');
  }
  
  static warn(message, data = null) {
    console.log(`[경고] ${message}`, data ? JSON.stringify(data) : '');
  }
  
  static error(message, data = null) {
    console.log(`[오류] ${message}`, data ? JSON.stringify(data) : '');
  }
  
  static debug(message, data = null) {
    if (__ENV.DEBUG === 'true') {
      console.log(`[디버그] ${message}`, data ? JSON.stringify(data) : '');
    }
  }
  
  static apiCall(method, url, status, duration) {
    const statusColor = status >= 400 ? '오류' : status >= 300 ? '경고' : '정보';
    console.log(`[${statusColor}] ${method} ${url} - ${status} (${duration}ms)`);
  }
}