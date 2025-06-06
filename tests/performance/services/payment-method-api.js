/**
 * k6 성능 테스트용 결제수단 API 서비스
 */
import { config } from "../config/environment.js";
import { BaseApiService } from "./base-api.js";
import { Logger } from "../utils/logger.js";

export class PaymentMethodApiService extends BaseApiService {
  constructor() {
    super(config.API.PAYMENT_METHODS);
  }

  /**
   * 필터를 사용하여 결제수단 목록 조회
   */
  getPaymentMethods(filters = {}) {
    const defaultFilters = {
      size: 100,
      page: 0
    };

    const queryParams = Object.assign({}, defaultFilters, filters);

    Logger.debug('필터로 결제수단 목록 조회', queryParams);

    const response = this.get('', queryParams);
    return this.getSuccessfulResponse(response, '결제수단 목록 조회');
  }

  /**
   * ID로 결제수단 조회
   */
  getPaymentMethodById(id) {
    Logger.debug('ID로 결제수단 조회', { id });

    const response = this.get(id.toString());
    return this.getSuccessfulResponse(response, `결제수단 ${id} 조회`);
  }

  /**
   * 새 결제수단 생성
   */
  createPaymentMethod(paymentMethodData) {
    Logger.debug('결제수단 생성', paymentMethodData);

    const response = this.post('', paymentMethodData);
    return this.getSuccessfulResponse(response, '결제수단 생성');
  }

  /**
   * 결제수단 수정
   */
  updatePaymentMethod(id, paymentMethodData) {
    Logger.debug('결제수단 수정', Object.assign({ id: id }, paymentMethodData));

    const response = this.put(id.toString(), paymentMethodData);
    return this.getSuccessfulResponse(response, `결제수단 ${id} 수정`);
  }

  /**
   * 결제수단 삭제
   */
  deletePaymentMethod(id) {
    Logger.debug('결제수단 삭제', { id });

    const response = this.delete(id.toString());
    return response.status >= 200 && response.status < 300;
  }

  /**
   * 목록에서 첫 번째 결제수단 ID 반환
   */
  getFirstPaymentMethodId(filters = {}) {
    const paymentMethods = this.getPaymentMethods(filters);

    if (paymentMethods && paymentMethods.content && paymentMethods.content.length > 0) {
      const firstPaymentMethod = paymentMethods.content[0];
      Logger.debug('첫 번째 결제수단 발견', { id: firstPaymentMethod.id });
      return firstPaymentMethod.id;
    }

    Logger.warn('목록에서 결제수단을 찾을 수 없음');
    return null;
  }

  /**
   * 목록에서 랜덤 결제수단 ID 반환
   */
  getRandomPaymentMethodId(filters = {}) {
    const paymentMethods = this.getPaymentMethods(filters);

    if (paymentMethods && paymentMethods.content && paymentMethods.content.length > 0) {
      const randomIndex = Math.floor(Math.random() * paymentMethods.content.length);
      const randomPaymentMethod = paymentMethods.content[randomIndex];
      Logger.debug('랜덤 결제수단 발견', {
        id: randomPaymentMethod.id,
        index: randomIndex,
        total: paymentMethods.content.length
      });
      return randomPaymentMethod.id;
    }

    Logger.warn('목록에서 결제수단을 찾을 수 없음');
    return null;
  }
}

export const paymentMethodApi = new PaymentMethodApiService();
