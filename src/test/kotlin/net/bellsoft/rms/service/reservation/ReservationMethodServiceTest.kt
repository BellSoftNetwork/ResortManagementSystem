package net.bellsoft.rms.service.reservation

import io.kotest.assertions.throwables.shouldThrow
import io.kotest.core.spec.style.BehaviorSpec
import io.kotest.matchers.shouldBe
import net.bellsoft.rms.domain.reservation.method.ReservationMethodRepository
import net.bellsoft.rms.exception.DataNotFoundException
import net.bellsoft.rms.exception.DuplicateDataException
import net.bellsoft.rms.fixture.baseFixture
import net.bellsoft.rms.service.reservation.dto.ReservationMethodCreateDto
import net.bellsoft.rms.service.reservation.dto.ReservationMethodPatchDto
import net.bellsoft.rms.util.TestDatabaseSupport
import org.openapitools.jackson.nullable.JsonNullable
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.data.domain.PageRequest
import org.springframework.data.domain.Sort
import org.springframework.data.repository.findByIdOrNull
import org.springframework.test.context.ActiveProfiles

@SpringBootTest
@ActiveProfiles("test")
internal class ReservationMethodServiceTest(
    private val testDatabaseSupport: TestDatabaseSupport,
    private val reservationMethodService: ReservationMethodService,
    private val reservationMethodRepository: ReservationMethodRepository,
) : BehaviorSpec(
    {
        val fixture = baseFixture

        Given("예약 수단이 없는 상황에서") {
            When("예약 수단을 생성하면") {
                val reservationMethod = reservationMethodService.create(
                    ReservationMethodCreateDto(
                        name = "네이버",
                        commissionRate = 0.1,
                    ),
                )

                Then("정상적으로 생성된다") {
                    reservationMethod.run {
                        name shouldBe "네이버"
                        commissionRate shouldBe 0.1
                    }
                }
            }

            When("존재하지 않는 예약 수단 정보를 조회하면") {
                val exception = shouldThrow<DataNotFoundException> {
                    reservationMethodService.find(-1)
                }

                Then("조회 불가 예외가 발생한다") {
                    exception.message shouldBe "존재하지 않는 예약 수단"
                }
            }

            When("존재하지 않는 예약 수단 정보 수정을 시도하면") {
                val exception = shouldThrow<DataNotFoundException> {
                    reservationMethodService.update(
                        -1,
                        ReservationMethodPatchDto(
                            name = JsonNullable.undefined(),
                            commissionRate = JsonNullable.undefined(),
                            requireUnpaidAmountCheck = JsonNullable.undefined(),
                        ),
                    )
                }

                Then("조회 불가 예외가 발생한다") {
                    exception.message shouldBe "존재하지 않는 예약 수단"
                }
            }

            When("존재하지 않는 예약 수단 정보 삭제를 시도하면") {
                val exception = shouldThrow<DataNotFoundException> {
                    reservationMethodService.delete(-1)
                }

                Then("조회 불가 예외가 발생한다") {
                    exception.message shouldBe "존재하지 않는 예약 수단"
                }
            }
        }

        Given("예약 수단이 등록된 상황에서") {
            val reservationMethod = reservationMethodRepository.save(fixture())

            When("동일한 예약 수단 생성을 시도하면") {
                val exception = shouldThrow<DuplicateDataException> {
                    reservationMethodService.create(
                        ReservationMethodCreateDto(
                            name = reservationMethod.name,
                            commissionRate = reservationMethod.commissionRate,
                        ),
                    )
                }

                Then("중복 생성으로 예외가 발생한다") {
                    exception.message shouldBe "이미 존재하는 예약 수단"
                }
            }

            When("예약 수단 리스트를 조회하면") {
                val reservationMethods = reservationMethodService.findAll(
                    PageRequest.of(0, 1, Sort.by("id").descending()),
                )

                Then("모든 예약 수단이 정상적으로 조회된다") {
                    reservationMethods.page.size shouldBe 1
                    reservationMethods.values.first().run {
                        name shouldBe reservationMethod.name
                        commissionRate shouldBe reservationMethod.commissionRate
                    }
                }
            }

            When("특정 예약 수단 정보를 조회하면") {
                val findReservationMethod = reservationMethodService.find(reservationMethod.id)

                Then("정상적으로 조회된다") {
                    findReservationMethod.run {
                        name shouldBe findReservationMethod.name
                        commissionRate shouldBe findReservationMethod.commissionRate
                    }
                }
            }

            When("등록한 예약 수단 정보 수정을 시도하면") {
                val updateReservationMethod = reservationMethodService.update(
                    reservationMethod.id,
                    ReservationMethodPatchDto(
                        name = JsonNullable.of("BSN"),
                        commissionRate = JsonNullable.of(0.2),
                        requireUnpaidAmountCheck = JsonNullable.undefined(),
                    ),
                )

                Then("정상적으로 수정된다") {
                    updateReservationMethod.run {
                        name shouldBe "BSN"
                        commissionRate shouldBe 0.2
                    }
                }
            }

            When("등록한 예약 수단 정보 삭제를 시도하면") {
                val isDeleted = reservationMethodService.delete(reservationMethod.id)

                Then("정상적으로 삭제된다") {
                    isDeleted shouldBe true
                    reservationMethodRepository.findByIdOrNull(reservationMethod.id) shouldBe null
                }
            }
        }

        afterSpec {
            testDatabaseSupport.clear()
        }
    },
)
