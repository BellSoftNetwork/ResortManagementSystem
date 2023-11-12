package net.bellsoft.rms.service.config

import com.ninjasquad.springmockk.MockkBean
import io.kotest.core.spec.style.BehaviorSpec
import io.kotest.matchers.shouldBe
import io.mockk.every
import net.bellsoft.rms.domain.user.UserRepository
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.test.context.ActiveProfiles

@SpringBootTest
@ActiveProfiles("test")
internal class ConfigServiceTest(
    private val configService: ConfigService,
    @MockkBean private val userRepository: UserRepository,
) : BehaviorSpec(
    {
        Given("가입한 사용자가 없는 상황에서") {
            every { userRepository.count() } returns 0

            When("앱 설정 정보 요청 시") {
                val appConfigDto = configService.getAppConfig()

                Then("회원 가입 가능 여부 플래그가 활성화된다.") {
                    appConfigDto.isAvailableRegistration shouldBe true
                }
            }
        }

        Given("기존에 가입한 사용자가 있는 상황에서") {
            every { userRepository.count() } returns 1

            When("앱 설정 정보 요청 시") {
                val appConfigDto = configService.getAppConfig()

                Then("회원 가입 가능 여부 플래그가 비활성화된다.") {
                    appConfigDto.isAvailableRegistration shouldBe false
                }
            }
        }
    },
)
