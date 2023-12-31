import axios from 'axios';

// Tùy chỉnh cấu hình Axios ở đây
const instance = axios.create({
  baseURL: 'http://localhost:8080', // Địa chỉ API cụ thể của bạn
  timeout: 5000,
});

const AxiosPlugin = {
  install: function (Vue) {
    Vue.prototype.$axios = instance;
  },
};

export default AxiosPlugin;
