function MyPromise(executor) {
  this.state = 'pending'; // 初始状态为待定
  this.value = undefined; // 存储结果值
  this.reason = undefined; // 存储拒绝原因
  this.onFulfilledCallbacks = []; // 存储成功回调
  this.onRejectedCallbacks = []; // 存储失败回调

  const resolve = (value) => {
    if (this.state === 'pending') {
      this.state = 'fulfilled';
      this.value = value;
      this.onFulfilledCallbacks.forEach((callback) => callback(value));
    }
  };

  const reject = (reason) => {
    if (this.state === 'pending') {
      this.state = 'rejected';
      this.reason = reason;
      this.onRejectedCallbacks.forEach((callback) => callback(reason));
    }
  };

  try {
    executor(resolve, reject);
  } catch (error) {
    reject(error);
  }
}

MyPromise.prototype.then = function (onFulfilled, onRejected) {
  if (this.state === 'fulfilled') {
    onFulfilled(this.value);
  } else if (this.state === 'rejected') {
    onRejected(this.reason);
  } else {
    this.onFulfilledCallbacks.push(onFulfilled);
    this.onRejectedCallbacks.push(onRejected);
  }
  return this; // 支持链式调用
};
