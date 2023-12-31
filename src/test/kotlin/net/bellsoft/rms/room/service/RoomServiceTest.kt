package net.bellsoft.rms.room.service

import io.kotest.assertions.assertSoftly
import io.kotest.assertions.throwables.shouldThrow
import io.kotest.core.spec.style.BehaviorSpec
import io.kotest.matchers.collections.shouldContainExactly
import io.kotest.matchers.shouldBe
import io.kotest.matchers.shouldNotBe
import net.bellsoft.rms.authentication.exception.UserNotFoundException
import net.bellsoft.rms.common.dto.service.EntityReferenceDto
import net.bellsoft.rms.common.exception.DataNotFoundException
import net.bellsoft.rms.common.exception.DuplicateDataException
import net.bellsoft.rms.common.util.SecurityTestSupport
import net.bellsoft.rms.common.util.TestDatabaseSupport
import net.bellsoft.rms.fixture.baseFixture
import net.bellsoft.rms.payment.repository.PaymentMethodRepository
import net.bellsoft.rms.reservation.entity.Reservation
import net.bellsoft.rms.reservation.repository.ReservationRepository
import net.bellsoft.rms.revision.type.HistoryType
import net.bellsoft.rms.room.dto.filter.RoomFilterDto
import net.bellsoft.rms.room.dto.service.RoomCreateDto
import net.bellsoft.rms.room.dto.service.RoomPatchDto
import net.bellsoft.rms.room.entity.Room
import net.bellsoft.rms.room.repository.RoomGroupRepository
import net.bellsoft.rms.room.repository.RoomRepository
import net.bellsoft.rms.room.type.RoomStatus
import net.bellsoft.rms.user.entity.User
import org.openapitools.jackson.nullable.JsonNullable
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.data.domain.PageRequest
import org.springframework.data.repository.findByIdOrNull
import org.springframework.test.context.ActiveProfiles
import java.time.LocalDate

@SpringBootTest
@ActiveProfiles("test")
internal class RoomServiceTest(
    private val testDatabaseSupport: TestDatabaseSupport,
    private val securityTestSupport: SecurityTestSupport,
    private val roomService: RoomService,
    private val roomRepository: RoomRepository,
    private val roomGroupRepository: RoomGroupRepository,
    private val reservationRepository: ReservationRepository,
    private val paymentMethodRepository: PaymentMethodRepository,
) : BehaviorSpec(
    {
        val fixture = baseFixture.new {
            property(Room::roomGroup) { roomGroupRepository.save(fixture()) }
        }
        val loginUser: User = fixture()

        beforeContainer {
            if (it.descriptor.isRootTest())
                securityTestSupport.login(loginUser)
        }

        Given("객실 정보가 없는 상황에서 로그인 후") {
            When("전체 객실 정보를 조회하면") {
                val entityListDto = roomService.findAll(
                    PageRequest.of(0, 10),
                    RoomFilterDto(),
                )

                Then("빈 객실 목록이 반환 된다") {
                    entityListDto.page.totalElements shouldBe 0
                }
            }

            When("존재하지 않는 객실 정보를 조회하면") {
                val exception = shouldThrow<DataNotFoundException> {
                    roomService.find(-1)
                }

                Then("예외가 발생한다") {
                    exception.message shouldBe "존재하지 않는 객실"
                }
            }

            When("신규 객실 정보를 등록하면") {
                val customFixture = fixture.new {
                    property(RoomCreateDto::roomGroup) { EntityReferenceDto.of(roomGroupRepository.save(fixture())) }
                }
                val roomCreateDto: RoomCreateDto = customFixture()
                val result = roomService.create(roomCreateDto)

                Then("등록된 객실 정보가 반환 된다") {
                    result.number shouldBe roomCreateDto.number
                }

                Then("생성 이력이 등록된다") {
                    val entityListDto = roomService.findHistory(result.id, PageRequest.of(0, 10))

                    entityListDto.page.totalElements shouldBe 1
                    entityListDto.values.first().let {
                        it.historyType shouldBe HistoryType.CREATED
                    }
                }

                Then("생성자 정보에 로그인된 계정 정보가 등록된다") {
                    val room = roomRepository.findByIdOrNull(result.id)!!

                    room.run {
                        createdBy.id shouldBe loginUser.id
                        updatedBy.id shouldBe loginUser.id
                    }
                }
            }

            When("존재하지 않는 객실 정보 수정을 시도하면") {
                val exception = shouldThrow<DataNotFoundException> {
                    roomService.update(-1, fixture())
                }

                Then("예외가 발생한다") {
                    exception.message shouldBe "존재하지 않는 객실"
                }
            }

            When("존재하지 않는 객실 삭제를 시도하면") {
                val exception = shouldThrow<DataNotFoundException> {
                    roomService.delete(-1)
                }

                Then("예외가 발생한다") {
                    exception.message shouldBe "존재하지 않는 객실"
                }
            }
        }

        Given("객실이 10개 등록된 상황에서") {
            val customFixture = fixture.new {
                property(RoomCreateDto::roomGroup) { EntityReferenceDto.of(roomGroupRepository.save(fixture())) }
            }
            val rooms = roomRepository.saveAll(fixture<List<Room>> { repeatCount { 10 } })

            When("전체 객실 정보를 조회하면") {
                val entityListDto = roomService.findAll(
                    PageRequest.of(0, 10),
                    RoomFilterDto(),
                )

                Then("10개의 객실 정보가 반환 된다") {
                    entityListDto.page.totalElements shouldBe 10
                }
            }

            When("존재하는 객실 정보를 조회하면") {
                val room = rooms[0]
                val result = roomService.find(room.id)

                Then("등록된 객실 정보가 반환 된다") {
                    result.id shouldBe room.id
                }
            }

            When("신규 객실 정보를 등록하면") {
                val roomCreateDto: RoomCreateDto = customFixture()
                val result = roomService.create(roomCreateDto)

                Then("등록된 객실 정보가 반환 된다") {
                    result.number shouldBe roomCreateDto.number
                }
            }

            When("동일한 객실 번호를 가진 객실 정보를 등록하면") {
                val exception = shouldThrow<DuplicateDataException> {
                    roomService.create(customFixture { property(RoomCreateDto::number) { rooms[0].number } })
                }

                Then("중복된 객실 번호로 등록에 실패한다") {
                    exception.message shouldBe "이미 존재하는 객실"
                }
            }

            When("존재하는 객실 정보 수정을 시도하면") {
                val newLoginUser = securityTestSupport.login()
                val room = rooms[0]
                val result = roomService.update(
                    room.id,
                    RoomPatchDto(
                        number = JsonNullable.of("UPDATED"),
                    ),
                )

                loginUser.id shouldNotBe newLoginUser.id

                Then("객실 정보가 정상적으로 수정된다") {
                    result.number shouldBe "UPDATED"
                }

                Then("수정 이력이 등록된다") {
                    val entityListDto = roomService.findHistory(room.id, PageRequest.of(0, 10))

                    entityListDto.page.totalElements shouldBe 2
                    entityListDto.values.toList().let {
                        it[0].historyType shouldBe HistoryType.CREATED
                        it[1].historyType shouldBe HistoryType.UPDATED
                        it[1].updatedFields shouldBe setOf("updatedBy", "number")
                        it[1].entity.number shouldBe "UPDATED"
                        it[1].entity.updatedBy.id shouldBe newLoginUser.id
                    }
                }

                Then("수정자 정보가 로그인된 계정 정보로 변경된다") {
                    roomRepository.findByIdOrNull(result.id)!!.run {
                        createdBy.id shouldBe loginUser.id
                        updatedBy.id shouldBe newLoginUser.id
                    }
                }
            }

            When("존재하는 객실 정보 삭제를 시도하면") {
                val newLoginUser = securityTestSupport.login()
                val room = rooms[0]

                loginUser.id shouldNotBe newLoginUser.id
                roomService.delete(room.id)

                Then("객실 정보가 정상적으로 삭제된다") {
                    roomRepository.existsById(room.id) shouldBe false
                }

                Then("삭제 이력이 등록된다") {
                    val entityListDto = roomService.findHistory(room.id, PageRequest.of(0, 10))

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

            When("로그아웃 후 존재하는 객실 정보 삭제를 시도하면") {
                securityTestSupport.logout()

                val room = rooms[0]
                val exception = shouldThrow<UserNotFoundException> { roomService.delete(room.id) }

                Then("객실 정보를 삭제할 수 없다") {
                    exception.message shouldBe "로그인 필요"
                    roomRepository.existsById(room.id) shouldBe true
                }
            }
        }

        Given("비활성 상태의 객실이 1개 등록된 상황에서") {
            roomRepository.save(fixture { property(Room::status) { RoomStatus.INACTIVE } })

            When("활성 상태의 예약 가능한 객실 정보를 조회하면") {
                val entityListDto = roomService.findAll(
                    PageRequest.of(0, 10),
                    RoomFilterDto(status = RoomStatus.NORMAL),
                )

                Then("0개의 객실 정보가 반환된다") {
                    entityListDto.page.totalElements shouldBe 0
                }
            }
        }

        Given("희망 기간 외 예약이 잡혀있어 예약이 가능한 객실이 4개있을 때") {
            val customFixture = fixture.new {
                property(Reservation::paymentMethod) { paymentMethodRepository.save(fixture()) }
            }

            val availableRooms = roomRepository.saveAll(
                listOf(
                    fixture {
                        property(Room::note) {
                            """[0]
                                기존 예약 기간: ##=
                                희망 예약 기간: =##
                            """.trimIndent()
                        }
                    },
                    fixture {
                        property(Room::note) {
                            """[1]
                                기존 예약 기간: =##
                                희망 예약 기간: ##=
                            """.trimIndent()
                        }
                    },
                    fixture {
                        property(Room::note) {
                            """[2]
                                기존 예약 기간: ##@@
                                희망 예약 기간: =##=
                            """.trimIndent()
                        }
                    },
                    fixture {
                        property(Room::note) {
                            """[3]
                                기존 예약 기간: ====
                                희망 예약 기간: =##=
                            """.trimIndent()
                        }
                    },
                ),
            )

            reservationRepository.saveAll(
                listOf(
                    customFixture<Reservation> {
                        property(Reservation::stayStartAt) { LocalDate.of(2023, 11, 9) }
                        property(Reservation::stayEndAt) { LocalDate.of(2023, 11, 10) }
                    }.apply { addRoom(availableRooms[0]) },
                    customFixture<Reservation> {
                        property(Reservation::stayStartAt) { LocalDate.of(2023, 11, 11) }
                        property(Reservation::stayEndAt) { LocalDate.of(2023, 11, 12) }
                    }.apply { addRoom(availableRooms[1]) },
                    customFixture<Reservation> {
                        property(Reservation::stayStartAt) { LocalDate.of(2023, 11, 9) }
                        property(Reservation::stayEndAt) { LocalDate.of(2023, 11, 10) }
                    }.apply { addRoom(availableRooms[2]) },
                    customFixture<Reservation> {
                        property(Reservation::stayStartAt) { LocalDate.of(2023, 11, 11) }
                        property(Reservation::stayEndAt) { LocalDate.of(2023, 11, 12) }
                    }.apply { addRoom(availableRooms[2]) },
                ),
            )

            When("기간 내 예약 가능한 객실 정보를 조회하면") {
                val entityListDto = roomService.findAll(
                    PageRequest.of(0, 10),
                    RoomFilterDto(
                        stayStartAt = LocalDate.of(2023, 11, 10),
                        stayEndAt = LocalDate.of(2023, 11, 11),
                    ),
                )

                Then("4개의 객실 정보가 반환 된다") {
                    assertSoftly {
                        entityListDto.page.totalElements shouldBe 4
                        entityListDto.values.map { it.note } shouldContainExactly availableRooms.map { it.note }
                    }
                }
            }
        }

        Given("희망 기간 내 연박 예약이 잡혀있어 예약이 불가능한 객실이 7개있을 때") {
            val customFixture = fixture.new {
                property(Reservation::paymentMethod) { paymentMethodRepository.save(fixture()) }
            }

            val reservedRooms = roomRepository.saveAll(
                listOf(
                    fixture {
                        property(Room::note) {
                            """[0]
                                기존 예약 기간: ###=
                                희망 예약 기간: =###
                            """.trimIndent()
                        }
                    },
                    fixture {
                        property(Room::note) {
                            """[1]
                                기존 예약 기간: ###
                                희망 예약 기간: ###
                            """.trimIndent()
                        }
                    },
                    fixture {
                        property(Room::note) {
                            """[2]
                                기존 예약 기간: ##=
                                희망 예약 기간: ###
                            """.trimIndent()
                        }
                    },
                    fixture {
                        property(Room::note) {
                            """[3]
                                기존 예약 기간: =###=
                                희망 예약 기간: #####
                            """.trimIndent()
                        }
                    },
                    fixture {
                        property(Room::note) {
                            """[4]
                                기존 예약 기간: =##
                                희망 예약 기간: ###
                            """.trimIndent()
                        }
                    },
                    fixture {
                        property(Room::note) {
                            """[5]
                                기존 예약 기간: =###
                                희망 예약 기간: ###=
                            """.trimIndent()
                        }
                    },
                    fixture {
                        property(Room::note) {
                            """[6]
                                기존 예약 기간: #####
                                희망 예약 기간: =###=
                            """.trimIndent()
                        }
                    },
                ),
            )

            val reservations = reservationRepository.saveAll(
                listOf(
                    customFixture<Reservation> {
                        property(Reservation::stayStartAt) { LocalDate.of(2023, 11, 9) }
                        property(Reservation::stayEndAt) { LocalDate.of(2023, 11, 11) }
                    }.apply { addRoom(reservedRooms[0]) },
                    customFixture<Reservation> {
                        property(Reservation::stayStartAt) { LocalDate.of(2023, 11, 10) }
                        property(Reservation::stayEndAt) { LocalDate.of(2023, 11, 20) }
                    }.apply { addRoom(reservedRooms[1]) },
                    customFixture<Reservation> {
                        property(Reservation::stayStartAt) { LocalDate.of(2023, 11, 10) }
                        property(Reservation::stayEndAt) { LocalDate.of(2023, 11, 11) }
                    }.apply { addRoom(reservedRooms[2]) },
                    customFixture<Reservation> {
                        property(Reservation::stayStartAt) { LocalDate.of(2023, 11, 11) }
                        property(Reservation::stayEndAt) { LocalDate.of(2023, 11, 19) }
                    }.apply { addRoom(reservedRooms[3]) },
                    customFixture<Reservation> {
                        property(Reservation::stayStartAt) { LocalDate.of(2023, 11, 19) }
                        property(Reservation::stayEndAt) { LocalDate.of(2023, 11, 20) }
                    }.apply { addRoom(reservedRooms[4]) },
                    customFixture<Reservation> {
                        property(Reservation::stayStartAt) { LocalDate.of(2023, 11, 19) }
                        property(Reservation::stayEndAt) { LocalDate.of(2023, 11, 21) }
                    }.apply { addRoom(reservedRooms[5]) },
                    customFixture<Reservation> {
                        property(Reservation::stayStartAt) { LocalDate.of(2023, 11, 1) }
                        property(Reservation::stayEndAt) { LocalDate.of(2023, 11, 30) }
                    }.apply { addRoom(reservedRooms[6]) },
                ),
            )

            When("기간 내 예약 가능한 객실 정보를 조회하면") {
                val entityListDto = roomService.findAll(
                    PageRequest.of(0, 10),
                    RoomFilterDto(
                        stayStartAt = LocalDate.of(2023, 11, 10),
                        stayEndAt = LocalDate.of(2023, 11, 20),
                    ),
                )

                Then("0개의 객실 정보가 반환 된다") {
                    assertSoftly {
                        entityListDto.page.totalElements shouldBe 0
                        entityListDto.values.map { it.note } shouldContainExactly emptyList()
                    }
                }
            }

            When("첫번째 예약을 수정할 경우 기간 내 예약 가능한 객실 정보를 조회하면") {
                val entityListDto = roomService.findAll(
                    PageRequest.of(0, 10),
                    RoomFilterDto(
                        stayStartAt = LocalDate.of(2023, 11, 10),
                        stayEndAt = LocalDate.of(2023, 11, 20),
                        excludeReservationId = reservations[0].id,
                    ),
                )

                Then("1개의 객실 정보가 반환 된다") {
                    entityListDto.page.totalElements shouldBe 1
                }
            }
        }

        Given("객실이 배정되지 않은 예약만 잡혀있어 희망 기간 내 예약 가능한 객실이 1개있을 때") {
            val customFixture = fixture.new {
                property(Reservation::paymentMethod) { paymentMethodRepository.save(fixture()) }
            }

            val reservedRoom = roomRepository.save(fixture())

            reservationRepository.save(
                customFixture {
                    property(Reservation::stayStartAt) { LocalDate.of(2023, 11, 1) }
                    property(Reservation::stayEndAt) { LocalDate.of(2023, 11, 30) }
                },
            )

            When("기간 내 예약 가능한 객실 정보를 조회하면") {
                val entityListDto = roomService.findAll(
                    PageRequest.of(0, 10),
                    RoomFilterDto(
                        stayStartAt = LocalDate.of(2023, 11, 10),
                        stayEndAt = LocalDate.of(2023, 11, 20),
                    ),
                )

                Then("1개의 객실 정보가 반환 된다") {
                    entityListDto.page.totalElements shouldBe 1
                    entityListDto.values.first().id shouldBe reservedRoom.id
                }
            }
        }

        afterSpec {
            testDatabaseSupport.clear()
        }
    },
)
