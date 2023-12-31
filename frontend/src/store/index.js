import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    userInfor: JSON.parse(localStorage.getItem('userInfor')) || {},
    isAdmin: JSON.parse(localStorage.getItem('isAdmin')) || false,
    isLogged: JSON.parse(localStorage.getItem('isLogged')) || false,
  },
  getters: {
    getUserInfor(state) {
      return state.userInfor
    },

    getIsAdmin(state) {
      return state.isAdmin
    },

    getIsLogged(state) {
      return state.isLogged
    }
  },
  mutations: {
    postUserInfor(state, userInfor) {
      state.userInfor = userInfor
      localStorage.setItem('userInfor', JSON.stringify(state.userInfor));
    },

    postIsAdmin(state, isAdmin) {
      state.isAdmin = isAdmin
      localStorage.setItem('isAdmin', JSON.stringify(state.isAdmin));
    },

    postIsLogged(state, isLogged) {
      state.isLogged = isLogged
      localStorage.setItem('isLogged', JSON.stringify(state.isLogged));
    },
  },
  actions: {
    postUserInfor({ commit }, userInfor) {
      commit('postUserInfor', userInfor);
    },

    postIsAdmin({ commit }, isAdmin) {
      commit('postIsAdmin', isAdmin);
    },

    postIsLogged({ commit }, isLogged) {
      commit('postIsLogged', isLogged);
    },
  },
  modules: {
  }
})
