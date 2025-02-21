// 定义队列
const macro: Function[] = [];
const micro: Function[] = [];
const timer: { time: number; callback: Function }[] = [];

// 模拟 queueMicrotask
function queueMicrotask(callback: VoidFunction) {
    micro.push(callback);
}

// 模拟 setTimeout
function setTimeout(callback: Function, timeout: number) {
    timer.push({time: Date.now() + timeout, callback});
}

// 处理微任务队列
function runMicrotasks() {
    while (micro.length > 0) {
        const task = micro.shift(); // 取出第一个微任务
        task && task(); // 执行微任务
    }
}

// 检查并更新 timer，将到期的任务放入宏任务队列
function checkTimers() {
    const now = Date.now();
    for (let i = 0; i < timer.length; i++) {
        if (timer[i].time <= now) {
            macro.push(timer[i].callback); // 到期的 timer 放入宏任务队列
            timer.splice(i, 1); // 从 timer 中移除
        }
    }
}

// 主 Event Loop
function eventLoop() {
    while (micro.length > 0 || macro.length > 0 || timer.length > 0) {
        // 1. 先执行所有微任务
        runMicrotasks();

        // 2. 检查 timer，将到期的任务移到宏任务队列
        checkTimers();

        // 3. 如果有微任务，优先执行
        if (micro.length > 0) {
            runMicrotasks();
        }
        // 4. 执行一个宏任务
        else if (macro.length > 0) {
            const macroTask = macro.shift(); // 取出第一个宏任务
            macroTask && macroTask(); // 执行宏任务
            checkTimers(); // 宏任务执行后再次检查 timer
        }
        // 5. 如果没有任务，短暂休眠（模拟真实环境）
        else {
            continue;
        }
    }
    console.log("Event Loop 结束，所有任务已执行");
}

// 测试代码
console.log("Start");
const t1 = Date.now()
setTimeout(() => {
    console.log("Timeout 1", Date.now() - t1);
    queueMicrotask(() => console.log("Microtask from Timeout 1"));
}, 1000);
const t2 = Date.now()
setTimeout(() => {
    console.log("Timeout 2", Date.now() - t2);
}, 1000);

queueMicrotask(() => {
    console.log("Microtask 1");
});

console.log("End");

// 启动 Event Loop
eventLoop();
