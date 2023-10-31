package net.bellsoft.rms.domain.room.event

import io.kotest.core.spec.style.BehaviorSpec
import io.kotest.matchers.shouldBe
import io.mockk.mockkStatic
import net.bellsoft.rms.domain.JpaEntityTest
import net.bellsoft.rms.domain.room.RoomRepository
import net.bellsoft.rms.domain.user.UserRepository
import net.bellsoft.rms.fixture.baseFixture
import org.springframework.data.repository.findByIdOrNull
import java.time.LocalDateTime

@JpaEntityTest
internal class RoomEventRepositoryTest(
    private val userRepository: UserRepository,
    private val roomRepository: RoomRepository,
    private val roomEventRepository: RoomEventRepository,
) : BehaviorSpec(
    {
        val fixture = baseFixture

        mockkStatic(LocalDateTime::class)

        Given("객실 이벤트가 생성된 상황에서") {
            val user = userRepository.save(fixture())
            val room = roomRepository.save(fixture())
            val roomEvent = roomEventRepository.save(
                fixture {
                    property(RoomEvent::user) { user }
                    property(RoomEvent::room) { room }
                },
            )

            When("등록된 id로 조회하면") {
                val selectedRoomEvent = roomEventRepository.findByIdOrNull(roomEvent.id)

                Then("정상적으로 조회된다") {
                    selectedRoomEvent?.id shouldBe roomEvent.id
                }
            }

            When("등록되지 않은 id로 조회하면") {
                val unknownRoomEvent = roomEventRepository.findByIdOrNull(-1)

                Then("빈 값이 조회된다") {
                    unknownRoomEvent shouldBe null
                }
            }

            When("객실 이벤트를 삭제하면") {
                val roomEventId = roomEvent.id
                roomEventRepository.delete(roomEvent)

                Then("객실 이벤트 정보가 조회되지 않는다") {
                    roomEventRepository.findByIdOrNull(roomEventId) shouldBe null
                }
            }
        }
    },
)
