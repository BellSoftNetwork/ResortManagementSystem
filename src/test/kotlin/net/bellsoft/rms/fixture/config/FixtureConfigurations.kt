package net.bellsoft.rms.fixture.config

import net.bellsoft.rms.fixture.domain.reservation.ReservationFixture
import net.bellsoft.rms.fixture.domain.reservation.event.ReservationEventFixture
import net.bellsoft.rms.fixture.domain.reservation.method.ReservationMethodFixture
import net.bellsoft.rms.fixture.domain.room.RoomFixture
import net.bellsoft.rms.fixture.domain.room.event.RoomEventFixture
import net.bellsoft.rms.fixture.domain.user.UserFixture

// NOTE: 신규 도메인 추가 시 해당 도메인 Configuration 생성 후 아래에 등록 필요
@Suppress("ktlint:experimental:property-naming")
private val domainConfigurations = listOf(
    UserFixture.BASE_CONFIGURATION,
    RoomFixture.BASE_CONFIGURATION,
    RoomEventFixture.BASE_CONFIGURATION,
    ReservationFixture.BASE_CONFIGURATION,
    ReservationEventFixture.BASE_CONFIGURATION,
    ReservationMethodFixture.BASE_CONFIGURATION,
)

// NOTE: 신규 DTO 추가 시 해당 DTO Configuration 생성 후 아래에 등록 필요
// @Suppress("ktlint:experimental:property-naming")
// private val dtoFixtureConfigurations = listOf()

// NOTE: 신규 설정 리스트 생성 시 해당 설정 리스트를 아래에 등록 필요
@Suppress("ktlint:experimental:property-naming")
val integratedFixtureConfigurations = domainConfigurations // + dtoFixtureConfigurations
