
#include "semaphore.hpp"

// @brief cpp_semaphore basic constructor
cpp_semaphore::cpp_semaphore()
{
	printk("Create semaphore %p\n", this);
	k_sem_init(&_sema_internal, 0, K_SEM_MAX_LIMIT);
}

/*
 * @brief wait for a semaphore
 *
 * Test a semaphore to see if it has been signaled.  If the signal
 * count is greater than zero, it is decremented.
 *
 * @return 1 when semaphore is available
 */
int cpp_semaphore::wait(void)
{
	k_sem_take(&_sema_internal, K_FOREVER);
	return 1;
}

/*
 * @brief wait for a semaphore within a specified timeout
 *
 * Test a semaphore to see if it has been signaled.  If the signal
 * count is greater than zero, it is decremented. The function
 * waits for timeout specified
 *
 * @param timeout the specified timeout in ticks
 *
 * @return 1 if semaphore is available, 0 if timed out
 */
int cpp_semaphore::wait(int timeout)
{
	return k_sem_take(&_sema_internal, K_MSEC(timeout));
}

// @brief This function signals the specified semaphore.
// @return N/A
void cpp_semaphore::give(void)
{
	k_sem_give(&_sema_internal);
}

