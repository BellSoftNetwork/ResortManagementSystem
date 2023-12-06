package net.bellsoft.rms.service.reservation

import io.kotest.assertions.assertSoftly
import io.kotest.assertions.throwables.shouldThrow
import io.kotest.core.spec.style.BehaviorSpec
import io.kotest.matchers.shouldBe
import io.kotest.matchers.shouldNotBe
import net.bellsoft.rms.component.history.type.HistoryType
import net.bellsoft.rms.domain.reservation.Reservation
import net.bellsoft.rms.domain.reservation.ReservationRepository
import net.bellsoft.rms.domain.reservation.ReservationRoom
import net.bellsoft.rms.domain.reservation.ReservationRoomRepository
import net.bellsoft.rms.domain.reservation.method.ReservationMethodRepository
import net.bellsoft.rms.domain.room.RoomRepository
import net.bellsoft.rms.domain.user.User
import net.bellsoft.rms.exception.DataNotFoundException
import net.bellsoft.rms.exception.UserNotFoundException
import net.bellsoft.rms.fixture.baseFixture
import net.bellsoft.rms.service.common.dto.EntityReferenceDto
import net.bellsoft.rms.service.reservation.dto.ReservationCreateDto
import net.bellsoft.rms.service.reservation.dto.ReservationFilterDto
import net.bellsoft.rms.service.reservation.dto.ReservationPatchDto
import net.bellsoft.rms.util.SecurityTestSupport
import net.bellsoft.rms.util.TestDatabaseSupport
import org.openapitools.jackson.nullable.JsonNullable
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.data.domain.PageRequest
import org.springframework.data.repository.findByIdOrNull
import org.springframework.test.context.ActiveProfiles
import java.time.LocalDate

@SpringBootTest
@ActiveProfiles("test")
internal class ReservationServiceTest(
    private val testDatabaseSupport: TestDatabaseSupport,
    private val securityTestSupport: SecurityTestSupport,
    private val reservationService: ReservationService,
    private val reservationRepository: ReservationRepository,
    private val reservationRoomRepository: ReservationRoomRepository,
    private val reservationMethodRepository: ReservationMethodRepository,
    private val roomRepository: RoomRepository,
) : BehaviorSpec(
    {
        val reservationMethod = reservationMethodRepository.save(baseFixture())
        val fixture = baseFixture.new {
            property(Reservation::reservationMethod) { reservationMethod }
            property(ReservationCreateDto::reservationMethodId) { reservationMethod.id }
        }
        val loginUser: User = fixture()

        beforeContainer {
            if (it.descriptor.isRootTest())
                securityTestSupport.login(loginUser)
        }

        Given("예약 정보가 없는 상황에서 로그인 후") {
            When("전체 예약 정보를 조회하면") {
                val entityListDto = reservationService.findAll(
                    PageRequest.of(0, 10),
                    ReservationFilterDto(),
                )

                Then("빈 예약 목록이 반환 된다") {
                    entityListDto.page.totalElements shouldBe 0
                }
            }

            When("존재하지 않는 예약 정보를 조회하면") {
                val exception = shouldThrow<DataNotFoundException> {
                    reservationService.findById(-1)
                }

                Then("예외가 발생한다") {
                    exception.message shouldBe "존재하지 않는 예약"
                }
            }

            When("신규 예약 정보를 등록하면") {
                val reservationCreateDto: ReservationCreateDto = fixture()
                val result = reservationService.create(reservationCreateDto)

                Then("등록된 예약 정보가 반환 된다") {
                    result.name shouldBe reservationCreateDto.name
                }

                Then("생성 이력이 등록된다") {
                    val entityListDto = reservationService.findHistory(result.id, PageRequest.of(0, 10))

                    entityListDto.page.totalElements shouldBe 1
                    entityListDto.values.first().let {
                        it.historyType shouldBe HistoryType.CREATED
                    }
                }

                Then("생성자 정보에 로그인된 계정 정보가 등록된다") {
                    val reservation = reservationRepository.findByIdOrNull(result.id)!!

                    reservation.run {
                        createdBy.id shouldBe loginUser.id
                        updatedBy.id shouldBe loginUser.id
                    }
                }
            }

            When("존재하지 않는 예약 정보 수정을 시도하면") {
                val exception = shouldThrow<DataNotFoundException> {
                    reservationService.update(-1, fixture())
                }

                Then("예외가 발생한다") {
                    exception.message shouldBe "존재하지 않는 예약"
                }
            }

            When("존재하지 않는 예약 삭제를 시도하면") {
                val exception = shouldThrow<DataNotFoundException> {
                    reservationService.delete(-1)
                }

                Then("예외가 발생한다") {
                    exception.message shouldBe "존재하지 않는 예약"
                }
            }
        }

        Given("객실이 배정된 예약이 등록된 상황에서") {
            val room = roomRepository.save(fixture())
            val reservation = reservationRepository.save(fixture())

            reservationRoomRepository.saveAll(listOf(ReservationRoom(reservation, room)))

            When("객실 배정을 해제를 요청하면") {
                val result = reservationService.update(
                    reservation.id,
                    ReservationPatchDto(
                        rooms = JsonNullable.of(setOf()),
                    ),
                )

                Then("정상적으로 객실 배정이 해제된다") {
                    result.rooms.size shouldBe 0
                }
            }

            When("객실 배정을 해제를 요청하지 않으면") {
                val result = reservationService.update(
                    reservation.id,
                    ReservationPatchDto(),
                )

                Then("객실 배정이 유지된다") {
                    result.rooms.size shouldBe 1
                    result.rooms.first().id shouldBe room.id
                }
            }
        }

        Given("예약이 10개 등록된 상황에서") {
            val reservations = reservationRepository.saveAll(fixture<List<Reservation>> { repeatCount { 10 } })

            When("전체 예약 정보를 조회하면") {
                val entityListDto = reservationService.findAll(
                    PageRequest.of(0, 10),
                    ReservationFilterDto(),
                )

                Then("10개의 예약 정보가 반환 된다") {
                    entityListDto.page.totalElements shouldBe 10
                }
            }

            When("존재하는 예약 정보를 조회하면") {
                val reservation = reservations[0]
                val result = reservationService.findById(reservation.id)

                Then("등록된 예약 정보가 반환 된다") {
                    result.id shouldBe reservation.id
                }
            }

            When("신규 예약 정보를 등록하면") {
                val reservationCreateDto: ReservationCreateDto = fixture()
                val result = reservationService.create(reservationCreateDto)

                Then("등록된 예약 정보가 반환 된다") {
                    result.name shouldBe reservationCreateDto.name
                }
            }

            When("다른 계정으로 로그인 후 존재하는 예약 정보 수정을 시도하면") {
                val newRoom = roomRepository.save(fixture())
                val newLoginUser = securityTestSupport.login()
                val reservation = reservations[0]
                val result = reservationService.update(
                    reservation.id,
                    ReservationPatchDto(
                        rooms = JsonNullable.of(setOf(EntityReferenceDto(newRoom.id))),
                        name = JsonNullable.of("UPDATED"),
                    ),
                )

                loginUser.id shouldNotBe newLoginUser.id

                Then("예약 정보가 정상적으로 수정된다") {
                    result.name shouldBe "UPDATED"
                    result.rooms.size shouldBe 1
                    result.rooms.first().id shouldBe newRoom.id
                }

                Then("수정 이력이 등록된다") {
                    val entityListDto = reservationService.findHistory(reservation.id, PageRequest.of(0, 10))

                    entityListDto.page.totalElements shouldBe 2
                    entityListDto.values.toList().let {
                        it[0].historyType shouldBe HistoryType.CREATED
                        it[1].historyType shouldBe HistoryType.UPDATED
                        it[1].updatedFields shouldBe setOf("updatedBy", "rooms", "name")
                        it[1].entity.name shouldBe "UPDATED"
                        it[1].entity.rooms.first().id shouldBe newRoom.id
                        it[1].entity.updatedBy.id shouldBe newLoginUser.id
                    }
                }

                Then("수정자 정보가 로그인된 계정 정보로 변경된다") {
                    reservationRepository.findByIdOrNull(result.id)!!.run {
                        createdBy.id shouldBe loginUser.id
                        updatedBy.id shouldBe newLoginUser.id
                    }
                }
            }

            When("존재하는 예약 정보 삭제를 시도하면") {
                val newLoginUser = securityTestSupport.login()
                val reservation = reservations[0]

                loginUser.id shouldNotBe newLoginUser.id
                reservationService.delete(reservation.id)

                Then("예약 정보가 정상적으로 삭제된다") {
                    reservationRepository.existsById(reservation.id) shouldBe false
                }

                Then("삭제 이력이 등록된다") {
                    val entityListDto = reservationService.findHistory(reservation.id, PageRequest.of(0, 10))

                    entityListDto.page.totalElements shouldBe 2
                    assertSoftly {
                        entityListDto.values.toList().let {
                            it[0].historyType shouldBe HistoryType.CREATED
                            it[0].entity.createdBy.id shouldBe loginUser.id
                            it[0].entity.updatedBy.id shouldBe loginUser.id

                            it[1].historyType shouldBe HistoryType.DELETED
                            it[1].entity.createdBy.id shouldBe loginUser.id
                            it[1].entity.updatedBy.id shouldBe newLoginUser.id
                        }
                    }
                }
            }

            When("로그아웃 후 존재하는 예약 정보 삭제를 시도하면") {
                securityTestSupport.logout()

                val reservation = reservations[0]
                val exception = shouldThrow<UserNotFoundException> { reservationService.delete(reservation.id) }

                Then("예약 정보를 삭제할 수 없다") {
                    exception.message shouldBe "로그인 필요"
                    reservationRepository.existsById(reservation.id) shouldBe true
                }
            }
        }

        Given("각각 입실일이 다른 예약이 4개 등록된 상황에서") {
            reservationRepository.saveAll(
                listOf(
                    fixture { property(Reservation::stayStartAt) { LocalDate.of(2023, 10, 31) } },
                    fixture { property(Reservation::stayStartAt) { LocalDate.of(2023, 11, 1) } },
                    fixture { property(Reservation::stayStartAt) { LocalDate.of(2023, 11, 30) } },
                    fixture { property(Reservation::stayStartAt) { LocalDate.of(2023, 12, 1) } },
                ),
            )

            When("11월 예약 정보를 조회하면") {
                val entityListDto = reservationService.findAll(
                    PageRequest.of(0, 10),
                    ReservationFilterDto(
                        stayStartAt = LocalDate.of(2023, 11, 1),
                        stayEndAt = LocalDate.of(2023, 11, 30),
                    ),
                )

                Then("2개의 예약 정보가 반환 된다") {
                    entityListDto.page.totalElements shouldBe 2
                }
            }
        }

        afterSpec {
            testDatabaseSupport.clear()
        }
    },
)
