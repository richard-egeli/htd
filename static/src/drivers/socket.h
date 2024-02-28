
/**
 * @file socket.h
 * @brief Header file for socket related functions and definitions.
 *
 * It provides an interface for creating, connecting, sending, and receiving data over sockets.
 */
#ifndef SOCKET_H_
#define SOCKET_H_

#include <stdint.h>

#include "zephyr/net/net_ip.h"
#include "zephyr/sys_clock.h"

/**
 * @brief The size of each transmission segment for socket communication.
 *
 * It specifies the maximum number of bytes that can be sent in a single
 * transmission.
 * The actual size of the transmission segment is 748, but there's a 40 byte header
 * that is added to each transmission, so the max we can send is 708 bytes without fragmentation
 *
 */
#define SOCKET_TX_SEG_SIZE 708

/**
 * @brief The size of the receive segment for the socket.
 * It specifies the maximum number of bytes that can be received in a single
 * transmission.
 * The actual size of the transmission segment is 1488, but there's a 40 byte header
 * that is added to each transmission, so the max we receive is 1448 bytes without fragmentation
 */
#define SOCKET_RX_SEG_SIZE 1448

/**
 * @brief Typedef for the socket listen callback function.
 *
 * This callback function is used to handle incoming data on a socket.
 * It takes a buffer and its length as parameters.
 *
 * @param buffer The buffer containing the incoming data.
 * @param buf_len The length of the buffer.
 */
typedef void (*socket_listen_callback_t)(const uint8_t* buffer, int16_t buf_len);

/**
 * @brief Function pointer type for socket timeout check.
 *
 * This function pointer type is used to define a callback function that checks
 * whether a socket timeout has occurred.
 *
 * @param timeout The timeout value to check against.
 * @return true if the timeout has occurred, false otherwise.
 */
typedef bool (*socket_timeout_ptr)(k_timeout_t);

/**
 *@brief The socket error type.
 *
 */
typedef enum socket_err {
    SOCKET_OK,
    ERR_SOCK_CREATE,
    ERR_SOCK_INVALID,
    ERR_SOCK_TCP_ONLY,
    ERR_SOCK_CONNECT,
    ERR_SOCK_CLOSE,
    ERR_SOCK_WRITE,
    ERR_SOCK_READ,
    ERR_SOCK_OVERFLOW,
    ERR_SOCK_TIMEOUT,
    ERR_SOCK_BIND,
    ERR_SOCK_LISTEN,
    ERR_SOCK_ACCEPT,
} socket_err_t;

typedef struct socket_info {
    int32_t handle;
    struct sockaddr_in remote_addr;
} socket_info_t;

/**
 * Creates a socket and initializes the socket_info structure with the specified host and port.
 *
 * @param host The host address to connect to.
 * @param port The port number to connect to.
 * @param socket_info Pointer to the socket_info_t structure to be initialized.
 * @return Returns an 8-bit integer indicating the success or failure of the socket creation.
 *         0 indicates success, while a negative value indicates an error.
 */
int8_t socket_create(const char* host, uint16_t port, socket_info_t* socket_info);

/**
 *@brief Create a socket and connect to the host.
 * A wrapper around the POSIX socket() function.
 * @param host The host to connect to.
 * @param port The port to connect to.
 * @param socket_handle The socket handle to use.
 * @return int8_t
 */
int8_t socket_connect(const socket_info_t* socket_info);

/**
 *@brief Close socket.
 * A wrapper around the POSIX close() function.
 * @param sock_in The socket information to use.
 * @return int8_t
 */
int8_t socket_close(const socket_info_t* sock_in);

/**
 *@brief Writes bytes to the socket. A wrapper around the POSIX send() function.
 *
 * @param sock_in The socket information to use.
 * @param buffer The bytes to write.
 * @return int8_t
 */
int8_t socket_write(const socket_info_t* sock_in, const uint8_t* buffer, uint16_t len);

/**
 *@brief Configures how long the socket will wait for data before returning.
 *
 * @param sock_in The socket information to use.
 * @param timeout_ms The timeout in milliseconds.
 * @return int8_t
 */
int8_t socket_timeout(const socket_info_t* sock_in, uint32_t timeout_ms);

/**
 *@brief Reads bytes from the socket. A wrapper around the POSIX recv() function.
 *
 * @param sock_in The socket information to use.
 * @param buffer The buffer to read into.
 * @param buf_len The length of the buffer.
 * @return int16_t The number of bytes read or a negative error code.
 */
int16_t socket_read(const socket_info_t* sock_in, uint8_t* buffer, int16_t buf_len);

/**
 * @brief Listen for incoming connections on a specified port.
 *
 * @param type The socket type.
 * @param port The port number to listen on.
 * @param buffer A pointer to a buffer to store incoming data.
 * @param buf_len The length of the buffer.
 * @param timeout_ms The timeout in milliseconds for waiting for incoming connections.
 * @return int16_t The amount of bytes read on success, or a negative value on failure.
 */
int16_t socket_listener(uint16_t port,
                        socket_timeout_ptr timeout,
                        socket_listen_callback_t callback);

/**
 * Returns a string representation of the given socket error code.
 *
 * @param err The socket error code to convert to a string.
 * @return A string representation of the given socket error code.
 */
const char* socket_error(int8_t err);

#endif  // SOCKET_H_
