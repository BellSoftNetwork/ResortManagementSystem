package net.bellsoft.rms.fixture.config

import net.bellsoft.rms.authentication.fixture.DeviceInfoFixture
import net.bellsoft.rms.authentication.fixture.LoginAttemptFixture
import net.bellsoft.rms.payment.fixture.PaymentMethodFixture
import net.bellsoft.rms.reservation.fixture.ReservationCreateDtoFixture
import net.bellsoft.rms.reservation.fixture.ReservationFixture
import net.bellsoft.rms.reservation.fixture.ReservationRoomFixture
import net.bellsoft.rms.room.fixture.RoomCreateDtoFixture
import net.bellsoft.rms.room.fixture.RoomFixture
import net.bellsoft.rms.room.fixture.RoomGroupCreateDtoFixture
import net.bellsoft.rms.room.fixture.RoomGroupFixture
import net.bellsoft.rms.user.fixture.UserCreateDtoFixture
import net.bellsoft.rms.user.fixture.UserFixture
import net.bellsoft.rms.user.fixture.UserRegistrationRequestFixture

// NOTE: 신규 도메인 추가 시 해당 도메인 Configuration 생성 후 아래에 등록 필요
@Suppress("ktlint:experimental:property-naming")
private val domainConfigurations = listOf(
    UserFixture.BASE_CONFIGURATION,
    RoomFixture.BASE_CONFIGURATION,
    RoomGroupFixture.BASE_CONFIGURATION,
    ReservationFixture.BASE_CONFIGURATION,
    ReservationRoomFixture.BASE_CONFIGURATION,
    PaymentMethodFixture.BASE_CONFIGURATION,
    LoginAttemptFixture.BASE_CONFIGURATION,
)

// NOTE: 신규 DTO 추가 시 해당 DTO Configuration 생성 후 아래에 등록 필요
// @Suppress("ktlint:experimental:property-naming")
private val dtoFixtureConfigurations = listOf(
    UserRegistrationRequestFixture.BASE_CONFIGURATION,
    UserCreateDtoFixture.BASE_CONFIGURATION,
    ReservationCreateDtoFixture.BASE_CONFIGURATION,
    RoomCreateDtoFixture.BASE_CONFIGURATION,
    RoomGroupCreateDtoFixture.BASE_CONFIGURATION,
    DeviceInfoFixture.BASE_CONFIGURATION,
)

// NOTE: 신규 설정 리스트 생성 시 해당 설정 리스트를 아래에 등록 필요
@Suppress("ktlint:experimental:property-naming")
val integratedFixtureConfigurations = domainConfigurations + dtoFixtureConfigurations
