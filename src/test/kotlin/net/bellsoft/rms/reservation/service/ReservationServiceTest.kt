package net.bellsoft.rms.reservation.service

import io.kotest.assertions.assertSoftly
import io.kotest.assertions.throwables.shouldThrow
import io.kotest.core.spec.style.BehaviorSpec
import io.kotest.matchers.shouldBe
import io.kotest.matchers.shouldNotBe
import io.kotest.matchers.string.shouldContain
import net.bellsoft.rms.authentication.exception.UserNotFoundException
import net.bellsoft.rms.common.dto.service.EntityReferenceDto
import net.bellsoft.rms.common.exception.DataNotFoundException
import net.bellsoft.rms.common.util.SecurityTestSupport
import net.bellsoft.rms.common.util.TestDatabaseSupport
import net.bellsoft.rms.fixture.baseFixture
import net.bellsoft.rms.payment.repository.PaymentMethodRepository
import net.bellsoft.rms.reservation.dto.filter.ReservationFilterDto
import net.bellsoft.rms.reservation.dto.response.StatisticsPeriodType
import net.bellsoft.rms.reservation.dto.service.ReservationCreateDto
import net.bellsoft.rms.reservation.dto.service.ReservationPatchDto
import net.bellsoft.rms.reservation.entity.Reservation
import net.bellsoft.rms.reservation.entity.ReservationRoom
import net.bellsoft.rms.reservation.exception.UnavailableRoomException
import net.bellsoft.rms.reservation.repository.ReservationRepository
import net.bellsoft.rms.reservation.repository.ReservationRoomRepository
import net.bellsoft.rms.revision.type.HistoryType
import net.bellsoft.rms.room.entity.Room
import net.bellsoft.rms.room.repository.RoomGroupRepository
import net.bellsoft.rms.room.repository.RoomRepository
import net.bellsoft.rms.user.entity.User
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
    private val paymentMethodRepository: PaymentMethodRepository,
    private val roomRepository: RoomRepository,
    private val roomGroupRepository: RoomGroupRepository,
) : BehaviorSpec(
    {
        val paymentMethod = paymentMethodRepository.save(baseFixture())
        val fixture = baseFixture.new {
            property(Reservation::paymentMethod) { paymentMethod }
            property(ReservationCreateDto::paymentMethod) { EntityReferenceDto(paymentMethod.id) }
            property(Room::roomGroup) { roomGroupRepository.save(fixture()) }
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
            val reservation = reservationRepository.save(fixture())
            val room: Room = roomRepository.save(fixture())

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

            When("숙박 기간을 연장하면") {
                val extendStayEndAt = reservation.stayEndAt.plusDays(1)
                val result = reservationService.update(
                    reservation.id,
                    ReservationPatchDto(stayEndAt = JsonNullable.of(extendStayEndAt)),
                )

                Then("정상적으로 연장된다") {
                    result.stayEndAt shouldBe extendStayEndAt
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

        Given("객실이 배정되고 각각 입실일이 다른 예약이 4개 등록된 상황에서") {
            val room: Room = roomRepository.save(fixture())

            val reservations = reservationRepository.saveAll(
                listOf(
                    fixture {
                        property(Reservation::stayStartAt) { LocalDate.of(2023, 10, 31) }
                        property(Reservation::stayEndAt) { LocalDate.of(2023, 11, 1) }
                        property(Reservation::price) { 100000 }
                        property(Reservation::peopleCount) { 2 }
                    },
                    fixture {
                        property(Reservation::stayStartAt) { LocalDate.of(2023, 11, 1) }
                        property(Reservation::stayEndAt) { LocalDate.of(2023, 11, 5) }
                        property(Reservation::price) { 200000 }
                        property(Reservation::peopleCount) { 3 }
                    },
                    fixture {
                        property(Reservation::stayStartAt) { LocalDate.of(2023, 11, 30) }
                        property(Reservation::stayEndAt) { LocalDate.of(2023, 12, 5) }
                        property(Reservation::price) { 300000 }
                        property(Reservation::peopleCount) { 4 }
                    },
                    fixture {
                        property(Reservation::stayStartAt) { LocalDate.of(2023, 12, 1) }
                        property(Reservation::stayEndAt) { LocalDate.of(2023, 12, 5) }
                        property(Reservation::price) { 400000 }
                        property(Reservation::peopleCount) { 5 }
                    },
                ),
            )
            reservationRoomRepository.saveAll(
                listOf(
                    ReservationRoom(reservations[0], room),
                    ReservationRoom(reservations[1], room),
                    ReservationRoom(reservations[2], room),
                    ReservationRoom(reservations[3], roomRepository.save(fixture())),
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

                Then("3개의 예약 정보가 반환 된다") {
                    entityListDto.page.totalElements shouldBe 3
                }
            }

            When("이미 객실이 배정된 기간에 추가 예약을 생성하려고 하면") {
                val exception = shouldThrow<UnavailableRoomException> {
                    reservationService.create(
                        fixture {
                            property(ReservationCreateDto::stayStartAt) { LocalDate.of(2023, 11, 1) }
                            property(ReservationCreateDto::stayEndAt) { LocalDate.of(2023, 11, 2) }
                            property(ReservationCreateDto::rooms) { setOf(EntityReferenceDto(room.id)) }
                        },
                    )
                }

                Then("중복 객실 배정 예외가 발생하면서 등록되지 않는다") {
                    exception.message shouldContain "(${room.number})"
                }
            }

            When("객실이 배정된 기존 예약을 이미 객실이 배정되어 배정 불가능한 기간으로 수정 요청 시") {
                val exception = shouldThrow<UnavailableRoomException> {
                    reservationService.update(
                        reservations.first().id,
                        ReservationPatchDto(
                            stayStartAt = JsonNullable.of(LocalDate.of(2023, 11, 1)),
                            stayEndAt = JsonNullable.of(LocalDate.of(2023, 11, 2)),
                            rooms = JsonNullable.of(setOf(EntityReferenceDto(room.id))),
                        ),
                    )
                }

                Then("중복 객실 설정 예외가 발생하면서 수정되지 않는다") {
                    exception.message shouldContain "(${room.number})"
                }
            }

            When("기본 기간 타입(MONTHLY)으로 예약 통계를 조회하면") {
                val startDate = LocalDate.of(2023, 10, 1)
                val endDate = LocalDate.of(2023, 12, 31)
                val statistics = reservationService.getStatistics(startDate, endDate)

                Then("월별 통계 데이터가 정상적으로 조회된다") {
                    statistics.periodType shouldBe StatisticsPeriodType.MONTHLY
                    statistics.stats.size shouldBe 3

                    // 10월 데이터 확인
                    val octoberStats = statistics.stats.find { it.period == "2023-10" }
                    octoberStats shouldNotBe null
                    octoberStats?.totalSales shouldBe 100000L
                    octoberStats?.totalReservations shouldBe 1
                    octoberStats?.totalGuests shouldBe 2

                    // 11월 데이터 확인
                    val novemberStats = statistics.stats.find { it.period == "2023-11" }
                    novemberStats shouldNotBe null
                    novemberStats?.totalSales shouldBe 500000L // 200000 + 300000
                    novemberStats?.totalReservations shouldBe 2
                    novemberStats?.totalGuests shouldBe 7 // 3 + 4

                    // 12월 데이터 확인
                    val decemberStats = statistics.stats.find { it.period == "2023-12" }
                    decemberStats shouldNotBe null
                    decemberStats?.totalSales shouldBe 400000L
                    decemberStats?.totalReservations shouldBe 1
                    decemberStats?.totalGuests shouldBe 5
                }
            }

            When("일별 기간 타입으로 예약 통계를 조회하면") {
                val startDate = LocalDate.of(2023, 11, 1)
                val endDate = LocalDate.of(2023, 11, 5)
                val statistics = reservationService.getStatistics(startDate, endDate, StatisticsPeriodType.DAILY)

                Then("일별 통계 데이터가 정상적으로 조회된다") {
                    statistics.periodType shouldBe StatisticsPeriodType.DAILY
                    statistics.stats.size shouldBe 2

                    // 11월 1일 데이터 확인
                    val nov1Stats = statistics.stats.find { it.period == "2023-11-01" }
                    nov1Stats shouldNotBe null
                    nov1Stats?.totalSales shouldBe 200000L
                    nov1Stats?.totalReservations shouldBe 1
                    nov1Stats?.totalGuests shouldBe 3
                }
            }
        }

        afterSpec {
            testDatabaseSupport.clear()
        }
    },
)
