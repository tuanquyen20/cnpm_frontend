<template>
  <v-app id="inspire">
    <v-main class="grey lighten-3">
      <v-container fluid class="remove-padding custom-background">
        <v-img src="@/assets/picture.jpg" class="fill-height blur"></v-img>

        <v-card :loading="loading" class="float" width="700">
          <v-card-title>Sign In</v-card-title>
          <v-card-actions>
            <v-row>
              <v-col cols="12" class="pl-6 pr-6"
                ><v-text-field
                  v-model="email"
                  label="Email"
                  :rules="emailRules"
                ></v-text-field
              ></v-col>
              <v-col cols="12" class="pl-6 pr-6"
                ><v-text-field
                  v-model="password"
                  type="password"
                  label="Password"
                ></v-text-field
              ></v-col>
              <v-col cols="3">
                <v-btn
                  @click="login(email, password)"
                  depressed
                  color="#252525"
                  class="white--text ml-3"
                  >Sign In</v-btn
                >
              </v-col>
            </v-row>
          </v-card-actions>
        </v-card>
      </v-container>
    </v-main>
  </v-app>
</template>

<script>
export default {
  data() {
    return {
      email: "",
      emailRules: [(v) => /.+@.+\..+/.test(v) || "E-mail must be valid"],
      password: "",
      loading: false,
      dialog: true,
    };
  },

  methods: {
    login(email, password) {
      this.$axios.get(`/accounts/email=${email}`).then((res) => {
        var user = res.data;
        if (password === user.password) {
          if (user.type === "admin") {
            this.loading = true;
            this.$store.dispatch("postUserInfor", user)
            this.$store.dispatch("postIsAdmin", true)
            this.$store.dispatch("postIsLogged", true)
            setTimeout(() => {
              this.$router.push({ name: "orders" });
            }, 3000);
          } else {
            this.loading = true;
            this.$store.dispatch("postUserInfor", user)
            this.$store.dispatch("postIsLogged", true)
            setTimeout(() => {
              this.$router.push({ name: "home" });
            }, 3000);
          }
        }
      })
      .catch((error) => {
        console.log(error);
      })
    },
  },
};
</script>

<style lang="scss" scoped>
.remove-padding {
  padding: 0 !important;
}

.blur {
  opacity: 0.5;
}

.float {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}
</style>
