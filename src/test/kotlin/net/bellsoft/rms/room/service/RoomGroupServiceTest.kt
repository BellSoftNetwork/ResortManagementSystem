package net.bellsoft.rms.room.service

import io.kotest.assertions.throwables.shouldThrow
import io.kotest.core.spec.style.BehaviorSpec
import io.kotest.matchers.shouldBe
import io.kotest.matchers.shouldNotBe
import net.bellsoft.rms.authentication.exception.UserNotFoundException
import net.bellsoft.rms.common.exception.DataNotFoundException
import net.bellsoft.rms.common.exception.DuplicateDataException
import net.bellsoft.rms.common.exception.RelatedDataException
import net.bellsoft.rms.common.util.SecurityTestSupport
import net.bellsoft.rms.common.util.TestDatabaseSupport
import net.bellsoft.rms.fixture.baseFixture
import net.bellsoft.rms.payment.repository.PaymentMethodRepository
import net.bellsoft.rms.reservation.entity.Reservation
import net.bellsoft.rms.reservation.repository.ReservationRepository
import net.bellsoft.rms.room.dto.filter.RoomFilterDto
import net.bellsoft.rms.room.dto.service.RoomGroupCreateDto
import net.bellsoft.rms.room.dto.service.RoomGroupPatchDto
import net.bellsoft.rms.room.entity.Room
import net.bellsoft.rms.room.entity.RoomGroup
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
internal class RoomGroupServiceTest(
    private val testDatabaseSupport: TestDatabaseSupport,
    private val securityTestSupport: SecurityTestSupport,
    private val roomGroupService: RoomGroupService,
    private val roomRepository: RoomRepository,
    private val roomGroupRepository: RoomGroupRepository,
    private val reservationRepository: ReservationRepository,
    private val paymentMethodRepository: PaymentMethodRepository,
) : BehaviorSpec(
    {
        val fixture = baseFixture.new {
            property(Reservation::paymentMethod) { paymentMethodRepository.save(baseFixture()) }
        }
        val loginUser: User = fixture()

        beforeContainer {
            if (it.descriptor.isRootTest())
                securityTestSupport.login(loginUser)
        }

        Given("객실 그룹 정보가 없는 상황에서 로그인 후") {
            When("전체 객실 그룹 정보를 조회하면") {
                val entityListDto = roomGroupService.findAll(
                    PageRequest.of(0, 10),
                )

                Then("빈 객실 그룹 목록이 반환 된다") {
                    entityListDto.page.totalElements shouldBe 0
                }
            }

            When("존재하지 않는 객실 그룹 정보를 조회하면") {
                val exception = shouldThrow<DataNotFoundException> {
                    roomGroupService.find(
                        -1,
                        RoomFilterDto(),
                    )
                }

                Then("예외가 발생한다") {
                    exception.message shouldBe "존재하지 않는 객실 그룹"
                }
            }

            When("신규 객실 그룹 정보를 등록하면") {
                val roomGroupCreateDto: RoomGroupCreateDto = fixture()
                val result = roomGroupService.create(roomGroupCreateDto)

                Then("등록된 객실 정보가 반환 된다") {
                    result.name shouldBe roomGroupCreateDto.name
                }

                Then("생성자 정보에 로그인된 계정 정보가 등록된다") {
                    val roomGroup = roomGroupRepository.findByIdOrNull(result.id)!!

                    roomGroup.run {
                        createdBy.id shouldBe loginUser.id
                        updatedBy.id shouldBe loginUser.id
                    }
                }
            }

            When("존재하지 않는 객실 그룹 정보 수정을 시도하면") {
                val exception = shouldThrow<DataNotFoundException> {
                    roomGroupService.update(-1, fixture())
                }

                Then("예외가 발생한다") {
                    exception.message shouldBe "존재하지 않는 객실 그룹"
                }
            }

            When("존재하지 않는 객실 삭제를 시도하면") {
                val exception = shouldThrow<DataNotFoundException> {
                    roomGroupService.delete(-1)
                }

                Then("예외가 발생한다") {
                    exception.message shouldBe "존재하지 않는 객실 그룹"
                }
            }
        }

        Given("객실 그룹이 10개 등록된 상황에서") {
            val roomGroups = roomGroupRepository.saveAll(fixture<List<RoomGroup>> { repeatCount { 10 } })

            When("전체 객실 그룹 정보를 조회하면") {
                val entityListDto = roomGroupService.findAll(
                    PageRequest.of(0, 10),
                )

                Then("10개의 객실 그룹 정보가 반환 된다") {
                    entityListDto.page.totalElements shouldBe 10
                }
            }

            When("존재하는 객실 그룹 정보를 조회하면") {
                val roomGroup = roomGroups[0]
                val result = roomGroupService.find(
                    roomGroup.id,
                    RoomFilterDto(),
                )

                Then("등록된 객실 정보가 반환 된다") {
                    result.id shouldBe roomGroup.id
                }
            }

            When("신규 객실 그룹 정보를 등록하면") {
                val roomGroupCreateDto: RoomGroupCreateDto = fixture()
                val result = roomGroupService.create(roomGroupCreateDto)

                Then("등록된 객실 그룹 정보가 반환 된다") {
                    result.name shouldBe roomGroupCreateDto.name
                }
            }

            When("동일한 이름을 가진 객실 그룹 정보를 등록하면") {
                val exception = shouldThrow<DuplicateDataException> {
                    roomGroupService.create(fixture { property(RoomGroupCreateDto::name) { roomGroups[0].name } })
                }

                Then("중복된 객실 그룹명으로 등록에 실패한다") {
                    exception.message shouldBe "이미 존재하는 객실 그룹"
                }
            }

            When("존재하는 객실 그룹 정보 수정을 시도하면") {
                val newLoginUser = securityTestSupport.login()
                val roomGroup = roomGroups[0]
                val result = roomGroupService.update(
                    roomGroup.id,
                    RoomGroupPatchDto(
                        name = JsonNullable.of("UPDATED"),
                    ),
                )

                loginUser.id shouldNotBe newLoginUser.id

                Then("객실 그룹 정보가 정상적으로 수정된다") {
                    result.name shouldBe "UPDATED"
                }

                Then("수정자 정보가 로그인된 계정 정보로 변경된다") {
                    roomGroupRepository.findByIdOrNull(result.id)!!.run {
                        createdBy.id shouldBe loginUser.id
                        updatedBy.id shouldBe newLoginUser.id
                    }
                }
            }

            When("존재하는 객실 그룹 정보 삭제를 시도하면") {
                val newLoginUser = securityTestSupport.login()
                val roomGroup = roomGroups[0]

                loginUser.id shouldNotBe newLoginUser.id
                roomGroupService.delete(roomGroup.id)

                Then("객실 그룹 정보가 정상적으로 삭제된다") {
                    roomGroupRepository.existsById(roomGroup.id) shouldBe false
                }
            }

            When("로그아웃 후 존재하는 객실 그룹 정보 삭제를 시도하면") {
                securityTestSupport.logout()

                val roomGroup = roomGroups[0]
                val exception = shouldThrow<UserNotFoundException> { roomGroupService.delete(roomGroup.id) }

                Then("객실 그룹 정보를 삭제할 수 없다") {
                    exception.message shouldBe "로그인 필요"
                    roomGroupRepository.existsById(roomGroup.id) shouldBe true
                }
            }
        }

        Given("객실 그룹 내에 비활성 상태의 객실이 1개 등록된 상황에서") {
            val roomGroup = roomGroupRepository.save(fixture())
            roomRepository.save(
                fixture {
                    property(Room::roomGroup) { roomGroup }
                    property(Room::status) { RoomStatus.INACTIVE }
                },
            )

            When("활성 상태의 예약 가능한 객실 정보를 조회하면") {
                val entityDto = roomGroupService.find(
                    roomGroup.id,
                    RoomFilterDto(status = RoomStatus.NORMAL),
                )

                Then("0개의 객실 정보가 반환된다") {
                    entityDto.rooms.size shouldBe 0
                }
            }

            When("존재하는 객실 그룹 정보 삭제를 시도하면") {
                val exception = shouldThrow<RelatedDataException> { roomGroupService.delete(roomGroup.id) }

                Then("객실 그룹 정보를 삭제할 수 없다") {
                    exception.message shouldBe "그룹 내 객실이 존재하여 삭제 불가"
                    roomGroupRepository.existsById(roomGroup.id) shouldBe true
                }
            }
        }

        Given("희망 기간 외 예약이 잡혀있어 객실 그룹 내 예약이 가능한 객실이 4개있을 때") {
            val roomGroup = roomGroupRepository.save(fixture())
            val customFixture = fixture.new {
                property(Room::roomGroup) { roomGroup }
            }
            val availableRooms = roomRepository.saveAll(
                listOf(
                    customFixture {
                        property(Room::note) {
                            """[0]
                                기존 예약 기간: ##=
                                희망 예약 기간: =##
                            """.trimIndent()
                        }
                    },
                    customFixture {
                        property(Room::note) {
                            """[1]
                                기존 예약 기간: =##
                                희망 예약 기간: ##=
                            """.trimIndent()
                        }
                    },
                    customFixture {
                        property(Room::note) {
                            """[2]
                                기존 예약 기간: ##@@
                                희망 예약 기간: =##=
                            """.trimIndent()
                        }
                    },
                    customFixture {
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
                    fixture<Reservation> {
                        property(Reservation::stayStartAt) { LocalDate.of(2023, 11, 9) }
                        property(Reservation::stayEndAt) { LocalDate.of(2023, 11, 10) }
                    }.apply { addRoom(availableRooms[0]) },
                    fixture<Reservation> {
                        property(Reservation::stayStartAt) { LocalDate.of(2023, 11, 11) }
                        property(Reservation::stayEndAt) { LocalDate.of(2023, 11, 12) }
                    }.apply { addRoom(availableRooms[1]) },
                    fixture<Reservation> {
                        property(Reservation::stayStartAt) { LocalDate.of(2023, 11, 9) }
                        property(Reservation::stayEndAt) { LocalDate.of(2023, 11, 10) }
                    }.apply { addRoom(availableRooms[2]) },
                    fixture<Reservation> {
                        property(Reservation::stayStartAt) { LocalDate.of(2023, 11, 11) }
                        property(Reservation::stayEndAt) { LocalDate.of(2023, 11, 12) }
                    }.apply { addRoom(availableRooms[2]) },
                ),
            )

            When("기간 내 예약 가능한 객실 정보를 조회하면") {
                val entityDto = roomGroupService.find(
                    roomGroup.id,
                    RoomFilterDto(
                        stayStartAt = LocalDate.of(2023, 11, 10),
                        stayEndAt = LocalDate.of(2023, 11, 11),
                    ),
                )

                Then("입실일 기준 최근 퇴실일이 가장 먼 순으로 4개의 객실 정보가 반환된다") {
                    entityDto.rooms.map { it.room.id } shouldBe listOf(
                        availableRooms[1].id,
                        availableRooms[3].id,
                        availableRooms[0].id,
                        availableRooms[2].id,
                    )
                }
            }
        }

        Given("희망 기간 내 연박 예약이 잡혀있어 예약이 불가능한 객실이 객실 그룹 내 7개있을 때") {
            val roomGroup = roomGroupRepository.save(fixture())
            val customFixture = fixture.new {
                property(Room::roomGroup) { roomGroup }
            }
            val reservedRooms = roomRepository.saveAll(
                listOf(
                    customFixture {
                        property(Room::note) {
                            """[0]
                                기존 예약 기간: ###=
                                희망 예약 기간: =###
                            """.trimIndent()
                        }
                    },
                    customFixture {
                        property(Room::note) {
                            """[1]
                                기존 예약 기간: ###
                                희망 예약 기간: ###
                            """.trimIndent()
                        }
                    },
                    customFixture {
                        property(Room::note) {
                            """[2]
                                기존 예약 기간: ##=
                                희망 예약 기간: ###
                            """.trimIndent()
                        }
                    },
                    customFixture {
                        property(Room::note) {
                            """[3]
                                기존 예약 기간: =###=
                                희망 예약 기간: #####
                            """.trimIndent()
                        }
                    },
                    customFixture {
                        property(Room::note) {
                            """[4]
                                기존 예약 기간: =##
                                희망 예약 기간: ###
                            """.trimIndent()
                        }
                    },
                    customFixture {
                        property(Room::note) {
                            """[5]
                                기존 예약 기간: =###
                                희망 예약 기간: ###=
                            """.trimIndent()
                        }
                    },
                    customFixture {
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
                    fixture<Reservation> {
                        property(Reservation::stayStartAt) { LocalDate.of(2023, 11, 9) }
                        property(Reservation::stayEndAt) { LocalDate.of(2023, 11, 11) }
                    }.apply { addRoom(reservedRooms[0]) },
                    fixture<Reservation> {
                        property(Reservation::stayStartAt) { LocalDate.of(2023, 11, 10) }
                        property(Reservation::stayEndAt) { LocalDate.of(2023, 11, 20) }
                    }.apply { addRoom(reservedRooms[1]) },
                    fixture<Reservation> {
                        property(Reservation::stayStartAt) { LocalDate.of(2023, 11, 10) }
                        property(Reservation::stayEndAt) { LocalDate.of(2023, 11, 11) }
                    }.apply { addRoom(reservedRooms[2]) },
                    fixture<Reservation> {
                        property(Reservation::stayStartAt) { LocalDate.of(2023, 11, 11) }
                        property(Reservation::stayEndAt) { LocalDate.of(2023, 11, 19) }
                    }.apply { addRoom(reservedRooms[3]) },
                    fixture<Reservation> {
                        property(Reservation::stayStartAt) { LocalDate.of(2023, 11, 19) }
                        property(Reservation::stayEndAt) { LocalDate.of(2023, 11, 20) }
                    }.apply { addRoom(reservedRooms[4]) },
                    fixture<Reservation> {
                        property(Reservation::stayStartAt) { LocalDate.of(2023, 11, 19) }
                        property(Reservation::stayEndAt) { LocalDate.of(2023, 11, 21) }
                    }.apply { addRoom(reservedRooms[5]) },
                    fixture<Reservation> {
                        property(Reservation::stayStartAt) { LocalDate.of(2023, 11, 1) }
                        property(Reservation::stayEndAt) { LocalDate.of(2023, 11, 30) }
                    }.apply { addRoom(reservedRooms[6]) },
                ),
            )

            When("기간 내 예약 가능한 객실 정보를 조회하면") {
                val entityDto = roomGroupService.find(
                    roomGroup.id,
                    RoomFilterDto(
                        stayStartAt = LocalDate.of(2023, 11, 10),
                        stayEndAt = LocalDate.of(2023, 11, 20),
                    ),
                )

                Then("0개의 객실 정보가 반환 된다") {
                    entityDto.rooms.size shouldBe 0
                }
            }

            When("첫번째 예약을 수정할 경우 기간 내 예약 가능한 객실 정보를 조회하면") {
                val entityDto = roomGroupService.find(
                    roomGroup.id,
                    RoomFilterDto(
                        stayStartAt = LocalDate.of(2023, 11, 10),
                        stayEndAt = LocalDate.of(2023, 11, 20),
                        excludeReservationId = reservations[0].id,
                    ),
                )

                Then("1개의 객실 정보가 반환 된다") {
                    entityDto.rooms.size shouldBe 1
                }
            }
        }

        Given("객실 그룹 내 객실들이 골고루 예약을 받았을 때") {
            val roomGroup = roomGroupRepository.save(fixture())
            val rooms = roomRepository.saveAll(
                fixture<List<Room>> {
                    property(Room::roomGroup) { roomGroup }

                    repeatCount { 5 }
                },
            )

            fun createReservation(startDay: Int, endDay: Int, rooms: List<Room>) {
                reservationRepository.save(
                    fixture<Reservation> {
                        property(Reservation::stayStartAt) { LocalDate.of(2023, 12, startDay) }
                        property(Reservation::stayEndAt) { LocalDate.of(2023, 12, endDay) }
                    }.apply { rooms.forEach { addRoom(it) } },
                )
            }

            /*
            #: 숙박 중
            -: 공실

            날짜      -> 123456789012345678901234
            rooms[0] -> ###--####---######-##---
            rooms[1] -> ##--###---##-##---####--
            rooms[2] -> -##-##--###--##-##---###
            rooms[3] -> -----##-#########----##-
            rooms[4] -> ###########-------######
             */
            createReservation(1, 2, listOf(rooms[0], rooms[1]))
            createReservation(2, 3, listOf(rooms[0], rooms[2]))
            createReservation(5, 6, listOf(rooms[1], rooms[2]))
            createReservation(6, 7, listOf(rooms[0], rooms[1], rooms[3]))
            createReservation(8, 9, listOf(rooms[0]))
            createReservation(9, 11, listOf(rooms[2]))
            createReservation(11, 12, listOf(rooms[1]))
            createReservation(13, 18, listOf(rooms[0]))
            createReservation(14, 15, listOf(rooms[1], rooms[2]))
            createReservation(17, 18, listOf(rooms[2]))
            createReservation(19, 22, listOf(rooms[1]))
            createReservation(20, 21, listOf(rooms[0]))
            createReservation(22, 24, listOf(rooms[2]))
            createReservation(9, 17, listOf(rooms[3]))
            createReservation(22, 23, listOf(rooms[3]))
            createReservation(1, 11, listOf(rooms[4]))
            createReservation(19, 24, listOf(rooms[4]))

            When("08 ~ 09 일 동안 숙박 가능한 객실 정보를 조회하면") {
                val findRoomGroup = roomGroupService.find(
                    roomGroup.id,
                    RoomFilterDto(
                        stayStartAt = LocalDate.of(2023, 12, 8),
                        stayEndAt = LocalDate.of(2023, 12, 9),
                    ),
                )

                Then("희망 입실일 기준 최근 퇴실일이 가장 먼 순으로 객실 정보가 반환된다") {
                    findRoomGroup.rooms.map { it.room.id } shouldBe listOf(
                        rooms[2].id,
                        rooms[1].id,
                        rooms[3].id,
                    )
                }
            }

            When("18 ~ 19 일 동안 숙박 가능한 객실 정보를 조회하면") {
                val findRoomGroup = roomGroupService.find(
                    roomGroup.id,
                    RoomFilterDto(
                        stayStartAt = LocalDate.of(2023, 12, 18),
                        stayEndAt = LocalDate.of(2023, 12, 19),
                    ),
                )

                Then("희망 입실일 기준 최근 퇴실일이 가장 먼 순으로 객실 정보가 반환된다") {
                    findRoomGroup.rooms.map { it.room.id } shouldBe listOf(
                        rooms[4].id,
                        rooms[1].id,
                        rooms[3].id,
                        rooms[0].id,
                        rooms[2].id,
                    )
                }
            }

            When("예약이 없었던 전 달 01 ~ 02 일 동안 숙박 가능한 객실 정보를 조회하면") {
                val findRoomGroup = roomGroupService.find(
                    roomGroup.id,
                    RoomFilterDto(
                        stayStartAt = LocalDate.of(2023, 11, 1),
                        stayEndAt = LocalDate.of(2023, 11, 2),
                    ),
                )

                Then("객실 생성 순으로 객실 정보가 반환된다") {
                    findRoomGroup.rooms.map { it.room.id } shouldBe listOf(
                        rooms[0].id,
                        rooms[1].id,
                        rooms[2].id,
                        rooms[3].id,
                        rooms[4].id,
                    )
                }
            }
        }

        Given("객실 그룹 내 객실이 배정되지 않은 예약만 잡혀있어 희망 기간 내 예약 가능한 객실이 1개있을 때") {
            val roomGroup = roomGroupRepository.save(fixture())
            roomRepository.save(fixture { property(Room::roomGroup) { roomGroup } })

            reservationRepository.saveAll(
                listOf(
                    fixture {
                        property(Reservation::stayStartAt) { LocalDate.of(2023, 11, 1) }
                        property(Reservation::stayEndAt) { LocalDate.of(2023, 11, 30) }
                    },
                ),
            )

            When("기간 내 예약 가능한 객실 정보를 조회하면") {
                val entityDto = roomGroupService.find(
                    roomGroup.id,
                    RoomFilterDto(
                        stayStartAt = LocalDate.of(2023, 11, 10),
                        stayEndAt = LocalDate.of(2023, 11, 20),
                    ),
                )

                Then("1개의 객실 정보가 반환 된다") {
                    entityDto.rooms.size shouldBe 1
                }
            }
        }

        afterSpec {
            testDatabaseSupport.clear()
        }
    },
)
