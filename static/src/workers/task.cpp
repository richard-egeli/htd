#include "task.hpp"
#include "..\utilities\semaphore.hpp"
#include "Application.hpp"

namespace Hardware {
namespace TaskWorker {

k_tid_t threadID;

cpp_semaphore semTask;

const int kStackSize = 500;
const int kWorkingTime = 5000;

k_thread taskThreadConfig;
K_THREAD_STACK_DEFINE(task_stack, kStackSize);

void taskWorker(void)
{
	k_timer timer;

	k_timer_init(&timer, NULL, NULL);

	while (1) 
    {
		semTask.wait();                                          // wait for external enabling
		printk("Task: start running...\n");
		k_timer_start(&timer, K_MSEC(kWorkingTime), K_NO_WAIT);  // Simulate a delay and a "Got Fix" state */
		k_timer_status_sync(&timer);
		printk("Task: Done!\n");
		Application::Continue();                                // Continue execution in main Application
	}
}

} // namespace TaskWorker
} // namespace Hardware

using namespace Hardware;

void TaskWorker::Start()
{
	threadID = k_thread_create(
        &taskThreadConfig, task_stack, kStackSize, 
        (k_thread_entry_t) taskWorker, 
        NULL, NULL, NULL, K_PRIO_COOP(7), 0, K_NO_WAIT);
	k_thread_name_set(threadID, "taskWorker");

}

void TaskWorker::Enable() { semTask.give(); }

k_tid_t TaskWorker::getThreadId() { return threadID; }