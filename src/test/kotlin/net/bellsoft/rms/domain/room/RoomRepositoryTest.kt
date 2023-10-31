package net.bellsoft.rms.domain.room

import io.kotest.core.spec.style.BehaviorSpec
import io.kotest.matchers.shouldBe
import io.mockk.mockkStatic
import net.bellsoft.rms.domain.JpaEntityTest
import net.bellsoft.rms.fixture.baseFixture
import org.springframework.data.repository.findByIdOrNull
import java.time.LocalDateTime

@JpaEntityTest
internal class RoomRepositoryTest(
    private val roomRepository: RoomRepository,
) : BehaviorSpec(
    {
        val fixture = baseFixture

        mockkStatic(LocalDateTime::class)

        Given("한 객실이 생성된 상황에서") {
            val room = roomRepository.save(fixture())

            When("등록된 객실 번호로 조회하면") {
                val selectedRoom = roomRepository.findByNumber(room.number)

                Then("정상적으로 조회된다") {
                    selectedRoom?.id shouldBe room.id
                }
            }

            When("등록되지 않은 객실 번호로 조회하면") {
                val unknownRoom = roomRepository.findByNumber("NON-EXIST-ROOM")

                Then("빈 값이 조회된다") {
                    unknownRoom shouldBe null
                }
            }

            When("객실을 삭제하면") {
                val roomId = room.id
                roomRepository.delete(room)

                Then("객실 정보가 조회되지 않는다") {
                    roomRepository.findByIdOrNull(roomId) shouldBe null
                }
            }
        }
    },
)
