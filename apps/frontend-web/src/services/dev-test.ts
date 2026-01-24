import { api } from "boot/axios";

export interface DevTestResponse {
  message: string;
  data: any;
}

export interface ReservationGenerationOptions {
  startDate?: string;
  endDate?: string;
  regularReservations?: number;
  monthlyReservations?: number;
}

export interface GenerateTestDataRequest {
  type: "essential" | "reservation" | "all";
  reservationOptions?: ReservationGenerationOptions;
}

export async function generateTestData(request: GenerateTestDataRequest) {
  const result = await api.post<DevTestResponse>("/api/v1/dev/test-data", request);
  return result.data;
}
