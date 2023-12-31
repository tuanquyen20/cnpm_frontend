import Vue from "vue";
import VueRouter from "vue-router";
import store from "@/store"; // Nếu store nằm ở thư mục gốc của dự án

Vue.use(VueRouter);

const routes = [
  {
    path: "/",
    redirect: "/home",
  },
  {
    path: "/home",
    name: "home",
    meta: {
      public: true,
    },
    component: () =>
      import(/* webpackChunkName: "about" */ "../views/Home.vue"),
  },
  {
    path: "/log-in",
    name: "log-in",
    meta: {
      public: true,
    },
    component: () =>
      import(/* webpackChunkName: "about" */ "../views/Login.vue"),
  },
  {
    path: "/admin-dashboard",
    name: "admin-dashboard",
    component: () =>
      import(/* webpackChunkName: "about" */ "../views/AdminDashboard.vue"),
    meta: {
      requiresAdmin: true,
    },
    children: [
      {
        path: "/rooms",
        name: "rooms",
        meta: {
          requiresAdmin: true,
        },
        component: () =>
          import(/* webpackChunkName: "about" */ "../components/Rooms.vue"),
      },
    
      {
        path: "/orders",
        name: "orders",
        meta: {
          requiresAdmin: true,
        },
        component: () =>
          import(/* webpackChunkName: "about" */ "../components/Orders.vue"),
      },


      {
        path: "/blogs",
        name: "blogs",
        meta: {
          requiresAdmin: true,
        },
        component: () =>
          import(/* webpackChunkName: "about" */ "../components/Blogs.vue"),
      },

      {
        path: "/users",
        name: "users",
        meta: {},
        component: () =>
          import(/* webpackChunkName: "about" */ "../components/Users.vue"),
      },

      {
        path: "/accounts",
        name: "accounts",
        meta: {
          requiresAdmin: true,
        },
        component: () =>
          import(/* webpackChunkName: "about" */ "../components/Accounts.vue"),
      },
    ],
  },

  {
    path: "/blogs/:id",
    name: "blog-page",
    meta: {
      public: true,
    },
    component: () =>
      import(/* webpackChunkName: "about" */ "../views/BlogPage.vue"),
  },

  {
    path: "/error-page",
    name: "error-page",
    meta: {
      public: true,
    },
    component: () =>
      import(/* webpackChunkName: "about" */ "../views/ErrorPage.vue"),
  },
];

const router = new VueRouter({
  routes,
});

router.beforeEach((to, from, next) => {
  const { isLogged, isAdmin } = store.state;

  if (to.meta.requiresAdmin) {
    if (!isLogged) {
      next("/log-in");
    } else if (isAdmin) {
      next();
    } else {
      next("/error-page");
    }
  } else if (to.meta.requiresCustomer) {
    if (!isLogged) {
      next("/log-in");
    } else if (!isAdmin) {
      next();
    } else {
      next("/error-page");
    }
  } else if (to.meta.public) {
    next();
  } else {
    next("/error-page");
  }
});

export default router;
