package net.bellsoft.rms.authentication.service

import io.kotest.assertions.assertSoftly
import io.kotest.assertions.throwables.shouldNotThrowAny
import io.kotest.assertions.throwables.shouldThrow
import io.kotest.core.spec.style.BehaviorSpec
import io.kotest.matchers.shouldBe
import net.bellsoft.rms.authentication.dto.DeviceInfoDto
import net.bellsoft.rms.authentication.entity.LoginAttempt
import net.bellsoft.rms.authentication.exception.TooManyRequestsException
import net.bellsoft.rms.authentication.fixture.DeviceInfoFixture
import net.bellsoft.rms.authentication.repository.LoginAttemptRepository
import net.bellsoft.rms.common.util.TestDatabaseSupport
import net.bellsoft.rms.fixture.baseFixture
import net.bellsoft.rms.fixture.util.feature
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.test.context.ActiveProfiles
import java.time.LocalDateTime

@SpringBootTest
@ActiveProfiles("test")
class LoginAttemptServiceTest(
    private val loginAttemptService: LoginAttemptService,
    private val loginAttemptRepository: LoginAttemptRepository,
    private val testDatabaseSupport: TestDatabaseSupport,
) : BehaviorSpec(
    {
        val fixture = baseFixture.feature(DeviceInfoFixture.Feature.WITH_INTERNAL_IP)

        Given("로그인 시도 서비스") {
            val username = "testuser"
            val deviceInfoDto: DeviceInfoDto = fixture()

            When("로그인 시도를 기록하면") {
                loginAttemptService.recordLoginAttempt(username, deviceInfoDto, true)

                Then("데이터베이스에 기록된다") {
                    val attempts = loginAttemptRepository.findAll()

                    assertSoftly {
                        attempts.size shouldBe 1
                        attempts[0].username shouldBe username
                        attempts[0].ipAddress shouldBe deviceInfoDto.ipAddress
                        attempts[0].successful shouldBe true
                        attempts[0].deviceFingerprint shouldBe deviceInfoDto.deviceFingerprint
                    }
                }
            }

            When("시간이 지나 이전 실패 시도가 윈도우를 벗어나면") {
                // 과거의 실패한 로그인 시도 기록 (윈도우 밖)
                val pastAttempt = LoginAttempt(
                    username = username,
                    ipAddress = deviceInfoDto.ipAddress,
                    successful = false,
                    attemptAt = LocalDateTime.now().minusHours(1),
                    deviceFingerprint = deviceInfoDto.deviceFingerprint,
                )
                loginAttemptRepository.save(pastAttempt)

                Then("로그인 시도가 허용된다") {
                    shouldNotThrowAny {
                        loginAttemptService.checkLoginAttempts(username, deviceInfoDto)
                    }
                }
            }

            When("디바이스 핑거프린트 관련 테스트") {
                // 성공한 로그인 시도 기록
                val deviceOneInfoDto: DeviceInfoDto =
                    fixture.feature(DeviceInfoFixture.Feature.WITH_WINDOWS_DEVICE_FINGERPRINT)()
                val deviceTwoInfoDto: DeviceInfoDto =
                    fixture.feature(DeviceInfoFixture.Feature.WITH_ANDROID_DEVICE_FINGERPRINT)()
                val deviceThreeInfoDto: DeviceInfoDto =
                    fixture.feature(DeviceInfoFixture.Feature.WITHOUT_DEVICE_FINGERPRINT)()

                loginAttemptService.recordLoginAttempt(username, deviceOneInfoDto, true)

                Then("디바이스 핑거프린트가 변경되면 isDeviceChanged가 true를 반환한다") {
                    val isChanged = loginAttemptService.isDeviceChanged(username, deviceTwoInfoDto)
                    isChanged shouldBe true
                }

                Then("디바이스 핑거프린트가 동일하면 isDeviceChanged가 false를 반환한다") {
                    val isChanged = loginAttemptService.isDeviceChanged(username, deviceOneInfoDto)
                    isChanged shouldBe false
                }

                Then("디바이스 핑거프린트가 빈 값이면 isDeviceChanged가 true를 반환한다") {
                    val isChanged = loginAttemptService.isDeviceChanged(username, deviceThreeInfoDto)
                    isChanged shouldBe true
                }
            }

            When("null 디바이스 핑거프린트로 로그인 시도를 기록하면") {
                val deviceInfoDto: DeviceInfoDto = fixture {
                    property(DeviceInfoDto::deviceFingerprint) { "" }
                }
                loginAttemptService.recordLoginAttempt(username, deviceInfoDto, true)

                Then("데이터베이스에 null 디바이스 핑거프린트로 기록된다") {
                    val attempts = loginAttemptRepository.findAll()

                    assertSoftly {
                        attempts.size shouldBe 1
                        attempts[0].username shouldBe username
                        attempts[0].ipAddress shouldBe deviceInfoDto.ipAddress
                        attempts[0].successful shouldBe true
                        attempts[0].deviceFingerprint shouldBe ""
                    }
                }
            }

            When("이전에 성공한 로그인 기록이 없으면") {
                val newUsername = "newuser"

                Then("isDeviceChanged가 false를 반환한다") {
                    val isChanged = loginAttemptService.isDeviceChanged(newUsername, deviceInfoDto)
                    isChanged shouldBe false
                }
            }

            // 요청받은 로그인 시도 제한 검증 테스트 케이스
            When("동일 IP와 ID로 로그인을 4번 실패하고 5번째 시도 시") {
                val deviceInfoDto: DeviceInfoDto = fixture()

                // 4번의 실패한 로그인 시도 기록 (최대 5번)
                repeat(4) {
                    loginAttemptService.recordLoginAttempt(username, deviceInfoDto, false)
                }

                Then("5번째 시도에서는 로그인이 허용되어야 한다") {
                    shouldNotThrowAny {
                        loginAttemptService.checkLoginAttempts(username, deviceInfoDto)
                    }
                }
            }

            When("동일 IP와 ID로 로그인을 5번 실패하고 6번째 시도 시") {
                val deviceInfoDto: DeviceInfoDto = fixture()

                // 5번의 실패한 로그인 시도 기록 (최대 5번)
                repeat(5) {
                    loginAttemptService.recordLoginAttempt(username, deviceInfoDto, false)
                }

                Then("6번째 시도에서는 TooManyRequestsException이 발생해야 한다") {
                    shouldThrow<TooManyRequestsException> {
                        loginAttemptService.checkLoginAttempts(username, deviceInfoDto)
                    }.message shouldBe "너무 많은 로그인 시도가 있었습니다. 30 분 후에 다시 시도해주세요."
                }

                Then("정상적인 계정 정보라도 로그인 시도가 거부된다") {
                    // 정상 계정 정보로 로그인 시도
                    shouldThrow<TooManyRequestsException> {
                        loginAttemptService.checkLoginAttempts(username, deviceInfoDto)
                    }
                }
            }

            When("동일 IP와 다른 ID로 로그인을 14번 실패하고 15번째 다른 ID로 시도 시") {
                val deviceInfoDto: DeviceInfoDto = fixture()
                val normalUsername = "normaluser"

                // 14번의 실패한 로그인 시도 기록 (IP 기반 최대 15번 = 5 * 3)
                repeat(14) {
                    val deviceInfoDto: DeviceInfoDto = fixture()
                    loginAttemptService.recordLoginAttempt("test-$it", deviceInfoDto, false)
                }

                Then("15번째 정상적인 다른 ID로 로그인 시도가 허용된다") {
                    shouldNotThrowAny {
                        loginAttemptService.checkLoginAttempts(normalUsername, deviceInfoDto)
                    }
                }
            }

            When("동일 IP와 다른 ID로 로그인을 15번 실패하고 16번째 시도 시") {
                val normalUsername = "normaluser"

                // 15번의 실패한 로그인 시도 기록 (IP 기반 최대 15번 = 5 * 3)
                repeat(15) {
                    val deviceInfoDto: DeviceInfoDto = fixture()
                    loginAttemptService.recordLoginAttempt("test-$it", deviceInfoDto, false)
                }

                Then("16번째 어떤 ID든 로그인 시도에서는 TooManyRequestsException이 발생해야 한다") {
                    shouldThrow<TooManyRequestsException> {
                        loginAttemptService.checkLoginAttempts(normalUsername, deviceInfoDto)
                    }.message shouldBe "이 IP 주소에서 너무 많은 로그인 시도가 있었습니다. 30 분 후에 다시 시도해주세요."
                }

                Then("정상적인 계정 정보라도 로그인 시도가 거부된다") {
                    // 정상 계정 정보로 로그인 시도
                    shouldThrow<TooManyRequestsException> {
                        loginAttemptService.checkLoginAttempts(normalUsername, deviceInfoDto)
                    }
                }
            }

            When("1번 IP에서 ID를 연속 로그인 실패하여 429 상태가 되고 2번 IP로 같은 ID 로그인 시도 시") {
                val adminUsername = "admin"

                // 1번 IP에서 관리자 ID로 5번 실패하여 429 상태
                repeat(5) {
                    loginAttemptService.recordLoginAttempt(adminUsername, deviceInfoDto, false)
                }

                // 1번 IP에서는 로그인 시도가 거부되는지 확인
                shouldThrow<TooManyRequestsException> {
                    loginAttemptService.checkLoginAttempts(adminUsername, deviceInfoDto)
                }.message shouldBe "너무 많은 로그인 시도가 있었습니다. 30 분 후에 다시 시도해주세요."

                Then("2번 IP에서는 동일 ID로 정상 로그인이 가능해야 한다") {
                    val otherDeviceInfoDto: DeviceInfoDto = fixture.feature(DeviceInfoFixture.Feature.WITH_LOCAL_IP)()

                    // 2번 IP에서는 로그인 시도가 허용되는지 확인
                    shouldNotThrowAny {
                        loginAttemptService.checkLoginAttempts(adminUsername, otherDeviceInfoDto)
                    }
                }
            }

            When("악의적인 사용자가 관리자 계정에 대해 지속적인 인증 실패를 발생시킬 때") {
                val adminUsername = "admin"
                val hackerDeviceInfoDto: DeviceInfoDto = fixture.feature(DeviceInfoFixture.Feature.WITH_HACKER_IP)()
                val normalDeviceInfoDto: DeviceInfoDto = fixture()

                // 공격자 IP에서 관리자 ID로 5번 실패하여 429 상태
                repeat(5) {
                    val deviceInfoDto: DeviceInfoDto = fixture.feature(DeviceInfoFixture.Feature.WITH_HACKER_IP)()
                    loginAttemptService.recordLoginAttempt(adminUsername, deviceInfoDto, false)
                }

                // 공격자 IP에서는 로그인 시도가 거부되는지 확인
                shouldThrow<TooManyRequestsException> {
                    loginAttemptService.checkLoginAttempts(adminUsername, hackerDeviceInfoDto)
                }.message shouldBe "너무 많은 로그인 시도가 있었습니다. 30 분 후에 다시 시도해주세요."

                Then("실제 관리자는 본인 IP에서 정상적으로 로그인이 가능해야 한다") {
                    // 관리자 IP에서는 로그인 시도가 허용되는지 확인
                    shouldNotThrowAny {
                        loginAttemptService.checkLoginAttempts(adminUsername, normalDeviceInfoDto)
                    }
                }
            }

            When("동일 IP, 동일 ID로 4번 로그인 실패 후 성공하고 다시 1번 실패 후 시도할 때") {
                val deviceInfoDto: DeviceInfoDto = fixture()
                val testUsername = "resetuser"

                // 4번의 실패한 로그인 시도 기록
                repeat(4) {
                    loginAttemptService.recordLoginAttempt(testUsername, deviceInfoDto, false)
                }

                // 5번째에 성공한 로그인 기록
                loginAttemptService.recordLoginAttempt(testUsername, deviceInfoDto, true)

                // 성공 후 1번 실패한 로그인 기록
                loginAttemptService.recordLoginAttempt(testUsername, deviceInfoDto, false)

                Then("이전 실패 이력은 무시되고 1번의 실패만 카운트되어 로그인이 허용되어야 한다") {
                    shouldNotThrowAny {
                        loginAttemptService.checkLoginAttempts(testUsername, deviceInfoDto)
                    }
                }
            }
        }

        afterTest {
            testDatabaseSupport.clear()
        }
    },
)
