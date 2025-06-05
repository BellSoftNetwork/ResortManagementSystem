package net.bellsoft.rms.common.util

import com.fasterxml.jackson.databind.JsonNode
import com.fasterxml.jackson.databind.ObjectMapper
import io.kotest.matchers.shouldBe
import io.kotest.matchers.shouldNotBe
import io.kotest.matchers.string.shouldContain
import org.springframework.http.HttpStatus
import org.springframework.test.web.servlet.MockMvc
import org.springframework.test.web.servlet.MvcResult
import org.springframework.test.web.servlet.request.MockHttpServletRequestBuilder
import org.springframework.test.web.servlet.result.MockMvcResultHandlers

/**
 * Extension functions to integrate MockMvc with Kotest assertions
 */

/**
 * A DSL for writing MockMvc tests in a more Kotest-friendly way
 */
class KotestMockMvcDsl(private val mockMvc: MockMvc) {
    /**
     * Performs a request and executes assertions on the result
     */
    fun request(requestBuilder: MockHttpServletRequestBuilder, assertions: MvcResult.() -> Unit) {
        val result = mockMvc.perform(requestBuilder)
            .andDo(MockMvcResultHandlers.print())
            .andReturn()

        // Apply assertions to the result
        result.assertions()
    }
}

/**
 * Extension function to create a KotestMockMvcDsl instance
 */
fun MockMvc.kotestDsl(): KotestMockMvcDsl {
    return KotestMockMvcDsl(this)
}

/**
 * Performs a request and returns the result for Kotest-style assertions
 */
fun MockMvc.performAndReturn(requestBuilder: MockHttpServletRequestBuilder): MvcResult {
    return this.perform(requestBuilder)
        .andDo(MockMvcResultHandlers.print())
        .andReturn()
}

/**
 * Extension function to verify the status code with Kotest assertions
 */
fun MvcResult.statusShouldBe(expected: HttpStatus) {
    this.response.status shouldBe expected.value()
}

/**
 * Extension function to verify the status code is OK (200)
 */
fun MvcResult.shouldBeOk() = statusShouldBe(HttpStatus.OK)

/**
 * Extension function to verify the status code is CREATED (201)
 */
fun MvcResult.shouldBeCreated() = statusShouldBe(HttpStatus.CREATED)

/**
 * Extension function to verify the status code is NO_CONTENT (204)
 */
fun MvcResult.shouldBeNoContent() = statusShouldBe(HttpStatus.NO_CONTENT)

/**
 * Extension function to verify the status code is BAD_REQUEST (400)
 */
fun MvcResult.shouldBeBadRequest() = statusShouldBe(HttpStatus.BAD_REQUEST)

/**
 * Extension function to verify the status code is UNAUTHORIZED (401)
 */
fun MvcResult.shouldBeUnauthorized() = statusShouldBe(HttpStatus.UNAUTHORIZED)

/**
 * Extension function to verify the status code is FORBIDDEN (403)
 */
fun MvcResult.shouldBeForbidden() = statusShouldBe(HttpStatus.FORBIDDEN)

/**
 * Extension function to verify the status code is NOT_FOUND (404)
 */
fun MvcResult.shouldBeNotFound() = statusShouldBe(HttpStatus.NOT_FOUND)

/**
 * Extension function to verify the status code is TOO_MANY_REQUESTS (429)
 */
fun MvcResult.shouldBeTooManyRequests() = statusShouldBe(HttpStatus.TOO_MANY_REQUESTS)

/**
 * Extension function to verify the status code is INTERNAL_SERVER_ERROR (500)
 */
fun MvcResult.shouldBeInternalServerError() = statusShouldBe(HttpStatus.INTERNAL_SERVER_ERROR)

/**
 * Extension function to get the response content as a string
 */
fun MvcResult.getContentAsString(): String {
    return this.response.contentAsString
}

/**
 * Extension function to parse the response content as JSON
 */
fun MvcResult.getContentAsJson(objectMapper: ObjectMapper): JsonNode {
    val content = this.getContentAsString()

    return objectMapper.readTree(content)
}

/**
 * Extension function to parse the response content as the specified type
 */
inline fun <reified T> MvcResult.getContentAs(objectMapper: ObjectMapper): T {
    val content = this.getContentAsString()

    return objectMapper.readValue(content, T::class.java)
}

/**
 * Extension function to verify that a JSON path exists
 */
fun MvcResult.jsonPathShouldExist(path: String, objectMapper: ObjectMapper) {
    val jsonNode = this.getContentAsJson(objectMapper)
    val node = jsonNode.at(convertToJsonPointer(path))

    node shouldNotBe null
    node.isMissingNode shouldBe false
}

/**
 * JSON 경로를 JsonPointer 형식으로 변환하는 함수
 */
private fun convertToJsonPointer(path: String): String {
    return if (path.startsWith("$.")) {
        "/" + path.substring(2).replace(".", "/")
    } else {
        path.replace("$", "/").replace(".", "/")
    }
}

/**
 * Extension function to verify that a JSON path has the expected value
 */
fun MvcResult.jsonPathShouldBe(path: String, expectedValue: String, objectMapper: ObjectMapper) {
    val jsonNode = this.getContentAsJson(objectMapper)
    val node = jsonNode.at(convertToJsonPointer(path))

    node.asText() shouldBe expectedValue
}

/**
 * Extension function to verify that a JSON path contains the expected value
 */
fun MvcResult.jsonPathShouldContain(path: String, expectedValue: String, objectMapper: ObjectMapper) {
    val jsonNode = this.getContentAsJson(objectMapper)
    val node = jsonNode.at(convertToJsonPointer(path))

    node.asText() shouldContain expectedValue
}

/**
 * Extension function to verify that a JSON path does not exist
 */
fun MvcResult.jsonPathShouldNotExist(path: String, objectMapper: ObjectMapper) {
    val jsonNode = this.getContentAsJson(objectMapper)
    val node = jsonNode.at(convertToJsonPointer(path))

    node.isMissingNode shouldBe true
}

/**
 * Extension function to verify that a JSON path is an array with the expected size
 */
fun MvcResult.jsonPathShouldBeArrayOfSize(path: String, expectedSize: Int, objectMapper: ObjectMapper) {
    val jsonNode = this.getContentAsJson(objectMapper)
    val node = jsonNode.at(convertToJsonPointer(path))

    node.isArray shouldBe true
    node.size() shouldBe expectedSize
}

/**
 * Extension function to extract a value from the response using the specified JSON path
 */
fun MvcResult.extractJsonPath(path: String, objectMapper: ObjectMapper): String {
    val jsonNode = this.getContentAsJson(objectMapper)

    return jsonNode.at(convertToJsonPointer(path)).asText()
}
