#include "socket.h"

#include "zephyr/kernel.h"
#include "zephyr/net/net_ip.h"
#include "zephyr/net/socket.h"

int8_t socket_timeout(const socket_info_t* sock_in, uint32_t timeout_ms) {
    struct timeval timeout = {
        .tv_sec  = (timeout_ms / 1000),         // seconds
        .tv_usec = (timeout_ms % 1000) * 1000,  // microseconds
    };

    // set the socket timeout
    if (setsockopt(sock_in->handle, SOL_SOCKET, SO_RCVTIMEO, &timeout, sizeof(timeout)) != 0) {
        return -ERR_SOCK_TIMEOUT;  // return a negative error if it failed
    };

    return SOCKET_OK;
}

int8_t socket_write(const socket_info_t* sock_in, const uint8_t* buffer, uint16_t len) {
    uint16_t offset = 0;  // the offset of the buffer
    uint16_t bytes;       // the number of bytes written

    do {
        bytes = sendto(sock_in->handle,
                       &buffer[offset],
                       len - offset,
                       0,
                       (struct sockaddr*)&sock_in->remote_addr,
                       sizeof(sock_in->remote_addr));

        if (bytes < 0) return -ERR_SOCK_WRITE;  // return a negative error if the write failed

        offset += bytes;  // increment the offset by the number of bytes written
    } while (offset < len);

    return SOCKET_OK;
}

int16_t socket_read(const socket_info_t* sock_in, uint8_t* buffer, int16_t buf_len) {
    int16_t bytes   = 0;  // the number of bytes read
    int16_t out_len = 0;  // the total number of bytes read

    do {
        socklen_t addr_len = sizeof(sock_in->remote_addr);  // the length of the address
        bytes              = recvfrom(sock_in->handle,
                         buffer + out_len,
                         buf_len - out_len,
                         0,
                         (struct sockaddr*)&sock_in->remote_addr,
                         &addr_len);  // read bytes from the socket

        if (bytes < 0) {
            if (out_len > 0) break;  // if we've read bytes, return them

            return -ERR_SOCK_READ;  // return a negative error if the read failed
        }

        out_len += bytes;  // increment the total number of bytes read
    } while (bytes != 0 && out_len < buf_len);

    return out_len;  // return the number of bytes read or a negative error code
}

int8_t socket_create(const char* host, uint16_t port, socket_info_t* socket_info) {
    socket_info->handle                 = socket(AF_INET, SOCK_STREAM, IPPROTO_TCP);
    socket_info->remote_addr.sin_family = AF_INET;      // use IPv4
    socket_info->remote_addr.sin_port   = htons(port);  // set the port

    // initialize the address with the host
    if (inet_pton(AF_INET, host, &socket_info->remote_addr.sin_addr) != 1) {
        return -ERR_SOCK_INVALID;  // return a negative error if the host is invalid
    }

    if (socket_info->handle < 0) {
        return -ERR_SOCK_CREATE;  // return a negative error if socket creation failed
    }

    return SOCKET_OK;
}

int8_t socket_close(const socket_info_t* sock_in) {
    if (close(sock_in->handle) != 0) {
        return -ERR_SOCK_CLOSE;
    }

    return SOCKET_OK;
}

int8_t socket_connect(const socket_info_t* socket_info) {
    int error = connect(socket_info->handle,
                        (struct sockaddr*)&socket_info->remote_addr,
                        sizeof(struct sockaddr_in));
    if (error < 0) {
        printk("ERROR: socket_connect failed %d\n", error);
        return -ERR_SOCK_CONNECT;  // return a negative error if it failed
    }

    return SOCKET_OK;
}

int16_t socket_listener(uint16_t port,
                        socket_timeout_ptr timeout_ptr,
                        socket_listen_callback_t callback) {
    socket_info_t sock_in;
    sock_in.handle                      = socket(AF_INET, SOCK_STREAM, IPPROTO_TCP);
    sock_in.remote_addr.sin_family      = AF_INET;            // use IPv4
    sock_in.remote_addr.sin_port        = htons(port);        // set the port
    sock_in.remote_addr.sin_addr.s_addr = htonl(INADDR_ANY);  // accept connections from any address
    int8_t error                        = 0;                  // the error code

    if (sock_in.handle < 0) return -ERR_SOCK_CREATE;  // return if socket creation failed

    error = bind(sock_in.handle,
                 (struct sockaddr*)&sock_in.remote_addr,
                 sizeof(sock_in.remote_addr));  // bind the socket to the address

    if (error < 0) {
        close(sock_in.handle);
        return -ERR_SOCK_BIND;  // return a negative error code if it failed
    }

    if (listen(sock_in.handle, 10) < 0) {
        close(sock_in.handle);
        return -ERR_SOCK_LISTEN;  // return a negative error code if it failed
    }

    while (timeout_ptr(K_NO_WAIT)) {
        socket_info_t client_sock_in = {0};
        socklen_t client_sock_len    = sizeof(client_sock_in.remote_addr);
        client_sock_in.handle        = accept(sock_in.handle,
                                       (struct sockaddr*)&client_sock_in.remote_addr,
                                       &client_sock_len);  // accept a connection

        if (client_sock_in.handle < 0) {
            if (errno == EAGAIN) continue;

            close(sock_in.handle);
            return -ERR_SOCK_ACCEPT;  // return a negative error code if it failed
        }

        uint8_t buffer[SOCKET_RX_SEG_SIZE] = {0};             // the buffer to read into
        uint16_t buf_len                   = sizeof(buffer);  // the length of the buffer
        int16_t length                     = socket_read(&client_sock_in, buffer, buf_len);

        close(client_sock_in.handle);  // close the client socket

        if (length <= 0) continue;  // if the read failed, try again

        callback(buffer, length);  // if we've read bytes, return them
    }

    close(sock_in.handle);  // close the server socket
    return 0;
}
