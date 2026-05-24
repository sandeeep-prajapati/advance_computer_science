# Chapter 17: Scheduling & Async

## @Scheduled — Scheduled Tasks

Spring's `@Scheduled` annotation runs methods on a fixed schedule. Enable it with `@EnableScheduling`.

```java
@SpringBootApplication
@EnableScheduling
public class MyApplication { ... }
```

---

## Scheduling Options

```java
@Component
public class ScheduledTasks {

    // Fixed delay: wait 5 seconds AFTER the previous run completes
    @Scheduled(fixedDelay = 5000)
    public void processQueue() {
        emailQueue.flush();
    }

    // Fixed rate: run every 10 seconds regardless of completion time
    @Scheduled(fixedRate = 10000)
    public void syncMetrics() {
        metricsService.sync();
    }

    // Cron expression
    @Scheduled(cron = "0 0 2 * * *")      // 2:00 AM daily
    public void dailyCleanup() {
        cleanupService.deleteExpiredTokens();
    }

    @Scheduled(cron = "0 0 9 * * MON-FRI")  // 9 AM weekdays
    public void sendDailySummary() {
        reportService.sendSummaryEmails();
    }

    // With initialDelay — wait before first run
    @Scheduled(fixedRate = 60000, initialDelay = 30000)
    public void heartbeat() {
        monitoringService.ping();
    }

    // External cron expression from config
    @Scheduled(cron = "${app.cleanup.schedule:0 0 2 * * *}")
    public void configurableCleanup() {
        cleanupService.run();
    }
}
```

### Cron Syntax: `second minute hour dayOfMonth month dayOfWeek`

```
"0 0 * * * *"          every hour
"*/15 * * * * *"       every 15 seconds
"0 0 8-17 * * MON-FRI" 8 AM to 5 PM, Monday to Friday
"0 0 12 1 * *"         noon on the 1st of every month
"0 30 9 * * ?"         9:30 AM daily (? = any day of week)
```

---

## Thread Pool for Schedulers

By default, all `@Scheduled` tasks share a **single thread** — a slow task blocks others.

```java
@Configuration
public class SchedulerConfig implements SchedulingConfigurer {

    @Override
    public void configureTasks(ScheduledTaskRegistrar taskRegistrar) {
        taskRegistrar.setScheduler(scheduledTaskExecutor());
    }

    @Bean(destroyMethod = "shutdown")
    public Executor scheduledTaskExecutor() {
        return Executors.newScheduledThreadPool(10);
    }
}
```

---

## @Async — Asynchronous Execution

`@Async` runs a method in a separate thread, returning immediately to the caller.

```java
@SpringBootApplication
@EnableAsync
public class MyApplication { ... }
```

```java
@Service
public class EmailService {

    // Fire and forget — returns void
    @Async
    public void sendWelcomeEmail(String email) {
        log.info("Sending email on thread: {}", Thread.currentThread().getName());
        // ... slow email sending logic
    }

    // Returns a Future
    @Async
    public CompletableFuture<EmailResult> sendEmail(EmailRequest request) {
        EmailResult result = mailSender.send(request);
        return CompletableFuture.completedFuture(result);
    }
}

@Service
public class UserService {
    private final EmailService emailService;

    public UserDto registerUser(CreateUserRequest request) {
        User user = userRepository.save(new User(request));
        emailService.sendWelcomeEmail(user.getEmail());  // returns immediately
        return toDto(user);
    }
}
```

**Important:** `@Async` only works when called from **outside** the bean. Self-calls (calling `@Async` method from the same class) won't work — Spring can't intercept them.

---

## Async Thread Pool Configuration

```java
@Configuration
@EnableAsync
public class AsyncConfig implements AsyncConfigurer {

    @Override
    public Executor getAsyncExecutor() {
        ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();
        executor.setCorePoolSize(5);
        executor.setMaxPoolSize(20);
        executor.setQueueCapacity(100);
        executor.setThreadNamePrefix("async-");
        executor.setRejectedExecutionHandler(new ThreadPoolExecutor.CallerRunsPolicy());
        executor.initialize();
        return executor;
    }

    @Override
    public AsyncUncaughtExceptionHandler getAsyncUncaughtExceptionHandler() {
        return (ex, method, params) ->
            log.error("Async exception in {}: {}", method.getName(), ex.getMessage(), ex);
    }
}
```

### Multiple Executor Pools:
```java
@Bean("emailExecutor")
public Executor emailExecutor() {
    ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();
    executor.setCorePoolSize(3);
    executor.setMaxPoolSize(10);
    executor.setThreadNamePrefix("email-");
    executor.initialize();
    return executor;
}

@Bean("reportExecutor")
public Executor reportExecutor() {
    ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();
    executor.setCorePoolSize(2);
    executor.setMaxPoolSize(5);
    executor.setThreadNamePrefix("report-");
    executor.initialize();
    return executor;
}

// Use specific executor
@Async("emailExecutor")
public void sendEmail(String email) { ... }

@Async("reportExecutor")
public CompletableFuture<Report> generateReport(Long id) { ... }
```

---

## Combining Async and Parallel Execution

```java
@Service
public class DashboardService {

    private final UserService userService;
    private final OrderService orderService;
    private final AnalyticsService analyticsService;

    public DashboardData loadDashboard(Long userId) throws ExecutionException, InterruptedException {
        // Run all three in parallel
        CompletableFuture<UserDto> userFuture = CompletableFuture.supplyAsync(
            () -> userService.findById(userId));

        CompletableFuture<List<OrderDto>> ordersFuture = CompletableFuture.supplyAsync(
            () -> orderService.getRecentOrders(userId, 10));

        CompletableFuture<AnalyticsDto> analyticsFuture = CompletableFuture.supplyAsync(
            () -> analyticsService.getSummary(userId));

        // Wait for all to complete
        CompletableFuture.allOf(userFuture, ordersFuture, analyticsFuture).join();

        return new DashboardData(
            userFuture.get(),
            ordersFuture.get(),
            analyticsFuture.get()
        );
    }
}
```

---

## Virtual Threads (Java 21+)

Java 21's virtual threads (Project Loom) allow massive concurrency with simple thread-per-request code:

```yaml
# application.yml
spring:
  threads:
    virtual:
      enabled: true   # Spring Boot 3.2+ — uses virtual threads for Tomcat and @Async
```

```java
// Or configure manually
@Bean
public AsyncTaskExecutor applicationTaskExecutor() {
    return new TaskExecutorAdapter(Executors.newVirtualThreadPerTaskExecutor());
}
```

With virtual threads:
- Each HTTP request gets its own virtual thread (millions can exist simultaneously)
- Blocking I/O (database, HTTP) yields the carrier thread — no thread pool starvation
- No need for reactive programming for I/O-bound workloads

---

## Programmatic Scheduling

```java
@Service
public class DynamicSchedulerService {

    private final TaskScheduler taskScheduler;
    private final Map<String, ScheduledFuture<?>> scheduledTasks = new ConcurrentHashMap<>();

    public void scheduleTask(String taskId, Runnable task, String cronExpression) {
        CronTrigger trigger = new CronTrigger(cronExpression);
        ScheduledFuture<?> future = taskScheduler.schedule(task, trigger);
        scheduledTasks.put(taskId, future);
        log.info("Scheduled task {} with cron: {}", taskId, cronExpression);
    }

    public void cancelTask(String taskId) {
        ScheduledFuture<?> future = scheduledTasks.remove(taskId);
        if (future != null) {
            future.cancel(false);
            log.info("Cancelled task: {}", taskId);
        }
    }
}
```

---

## Transaction + @Async

`@Async` methods run in a new thread — the original transaction does **not** propagate.

```java
@Service
public class OrderService {

    @Transactional
    public Order placeOrder(OrderRequest request) {
        Order order = orderRepository.save(new Order(request));
        // At this point transaction is open

        notificationService.sendConfirmation(order.getId());
        // This runs in a NEW thread with NO transaction
        // It will start its own transaction when it hits @Transactional methods

        return order;
    }
}
```

---

## Best Practices

| Practice | Why |
|---|---|
| Always configure a thread pool | Default `SimpleAsyncTaskExecutor` creates a new thread per call |
| Set `setQueueCapacity` | Prevents unbounded queue memory growth |
| Handle async exceptions | They're swallowed silently without `AsyncUncaughtExceptionHandler` |
| Use `CompletableFuture` over `Future` | Better composition, non-blocking API |
| Log thread names | Makes debugging async flows easier |
| Don't use `@Async` in the same class | Bypasses the proxy — won't work |
