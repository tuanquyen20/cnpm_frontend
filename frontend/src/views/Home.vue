<template>
  <v-app id="inspire">
    <v-app-bar app color="#252525">
      <v-row>
        <v-col md="2">
          <v-img
            :src="require('../assets/logo.svg')"
            class="my-3"
            contain
            height="25"
          />
        </v-col>

        <v-col md="2" offset-md="8" class="mt-3" v-if="!isLogged">
          <v-btn @click="openLoginPage" text class="white--text mb-2"
            >Sign In</v-btn
          >
        </v-col>

        <v-col class="mt-3" v-if="isLogged && isAdmin">
          <v-btn @click="goToMyOrders" text class="white--text mb-2"
            >Orders</v-btn
          >
        </v-col>

        <v-col md="2" offset-md="5" class="mt-3" v-if="isLogged && !isAdmin">
          <v-btn @click="logOut" text class="white--text mb-2">Sign Out</v-btn>
        </v-col>

        <v-col md="2" offset-md="5" class="mt-3" v-if="isLogged && isAdmin">
          <v-btn @click="logOut" text class="white--text mb-2">Sign Out</v-btn>
        </v-col>
      </v-row>
    </v-app-bar>

    <v-main class="grey lighten-3 remove-padding">
      <v-container fluid class="remove-padding">
        <v-row style="min-height: 100vh">
          <v-col cols="6" class="d-flex">
            <div class="align-self-center" style="margin-left: 130px">
              <div style="display: flex; font-size: 16px">
                <div class="mr-6">
                  <svg
                    width="22"
                    height="28"
                    viewBox="0 0 22 28"
                    fill="none"
                    xmlns="http://www.w3.org/2000/svg"
                    xmlns:xlink="http://www.w3.org/1999/xlink"
                  >
                    <rect
                      width="22"
                      height="28"
                      fill="url(#pattern0)"
                      fill-opacity="0.5"
                    />
                    <defs>
                      <pattern
                        id="pattern0"
                        patternContentUnits="objectBoundingBox"
                        width="1"
                        height="1"
                      >
                        <use
                          xlink:href="#image0_2_14"
                          transform="matrix(0.0111111 0 0 0.00873016 0 0.107143)"
                        />
                      </pattern>
                      <image
                        id="image0_2_14"
                        width="90"
                        height="90"
                        xlink:href="data:image/png;base64,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"
                      />
                    </defs>
                  </svg>
                </div>

                <div class="ml-6" style="opacity: 50%">Da Nang, Viet Nam</div>
              </div>

              <div class="mt-6" style="font-size: 40px; font-weight: bold">
                Meliã Danang
              </div>

              <div class="mt-6" style="font-size: 16px; opacity: 50%">
                19 Truong Sa, Hoa Hai, Ngu Hanh Son, Da Nang
              </div>

              <div class="mt-6" style="font-size: 16px; opacity: 50%">
                (+84) 36 666 6250
              </div>

              <div class="mt-6">
                <v-icon
                  @click="openFacebookFanpage"
                  x-large
                  class="mr-3 icon-hover"
                  >mdi-facebook</v-icon
                >

                <v-icon
                  @click="openInstaFanpage"
                  x-large
                  class="ml-3 icon-hover"
                  >mdi-instagram</v-icon
                >
              </div>
            </div>
          </v-col>

          <v-col cols="6">
            <v-img src="@/assets/picture.jpg" class="fill-height"></v-img>
          </v-col>
        </v-row>
      </v-container>

      <v-container fluid class="white--text" style="background-color: #252525">
        <div class="mt-6 mb-6" style="font-size: 25px; margin-left: 125px">
          Find your perfect room
        </div>

        <div style="margin-left: 125px; margin-right: 125px" class="mb-6">
          <v-carousel cycle hide-delimiter-background show-arrows-on-hover>
            <v-carousel-item contain v-for="(item, i) in slideItems" :key="i">
              <v-row class="fill-height">
                <v-col cols="4">
                  <v-card
                    color="#D9D9D9"
                    class="overflow-y-auto"
                    max-height="500"
                  >
                    <v-card-title>{{ item.name }}</v-card-title>
                    <v-card-subtitle>{{
                      formattedPrice(item.price) + " VNĐ/night"
                    }}</v-card-subtitle>
                    <v-card-text
                      style="text-align: justify"
                      class="mt-1 black--text"
                      >{{ item.description }}</v-card-text
                    >
                    <v-card-actions>
                      <v-btn
                        color="#252525"
                        class="white--text"
                        height="50px"
                        width="100%"
                        @click="showBookingForm(item)"
                        v-if="!isAdmin"
                      >
                        Book now
                      </v-btn>
                    </v-card-actions>
                  </v-card>
                </v-col>
                <v-col cols="8">
                  <v-img class="rounded-lg" :src="item.background_link"></v-img>
                </v-col>
              </v-row>
            </v-carousel-item>
          </v-carousel>
        </div>
      </v-container>

      <v-container fluid class="remove-padding">
        <HomeBlogs
          :isAdmin="isAdmin"
          :isLogged="isLogged"
          :userInfor="userInfor"
        ></HomeBlogs>
        <v-img src="@/assets/picture.jpg" class="fill-height"></v-img>
      </v-container>
    </v-main>

    <v-dialog v-model="bookingDialog" width="900">
      <v-card>
        <v-card-title class="text-h5 grey lighten-2">
          Book Meliã Karaoke Room
        </v-card-title>

        <v-card-actions>
          <v-row>
            <v-col cols="6">
              <v-text-field
                v-model="bookingForm.name"
                label="Name"
              ></v-text-field>
            </v-col>
            <v-col cols="6">
              <v-select
                v-model="bookingForm.select"
                :hint="`${formattedPrice(bookingForm.select.price)} VNĐ/night`"
                :items="slideItems"
                item-text="name"
                item-value="id"
                label="Room Type"
                persistent-hint
                return-object
                single-line
              ></v-select>
            </v-col>
            <v-col cols="6">
              <v-text-field
                v-model="bookingForm.phone_number"
                label="Phone number"
              ></v-text-field>
            </v-col>
            <v-col cols="6">
              <v-text-field
                v-model="bookingForm.email"
                label="Email"
              ></v-text-field>
            </v-col>
            <v-col cols="12">
              <v-dialog v-model="dateRangeDialog" width="290px">
                <template v-slot:activator="{ on, attrs }">
                  <v-text-field
                    v-model="bookingForm.date_range"
                    label="Check-in and check-out date"
                    append-icon="mdi-calendar"
                    readonly
                    v-bind="attrs"
                    v-on="on"
                  ></v-text-field>
                </template>
                <v-date-picker
                  v-if="dateRangeDialog"
                  v-model="bookingForm.date_range"
                  range
                  full-width
                  color="#252525"
                  :min="date"
                >
                </v-date-picker>
              </v-dialog>
            </v-col>

            <v-col cols="6">
              <v-dialog v-model="timeStartDialog" width="290px">
                <template v-slot:activator="{ on, attrs }">
                  <v-text-field
                    v-model="bookingForm.time_start"
                    label="Check-in time"
                    append-icon="mdi-clock-time-four-outline"
                    readonly
                    v-bind="attrs"
                    v-on="on"
                  ></v-text-field>
                </template>
                <v-time-picker
                  v-if="timeStartDialog"
                  v-model="bookingForm.time_start"
                  full-width
                  color="#252525"
                  :max="
                    !bookingForm.date_range[1] ||
                    bookingForm.date_range[0] === bookingForm.date_range[1]
                      ? bookingForm.time_end
                      : `23:59`
                  "
                >
                </v-time-picker>
              </v-dialog>
            </v-col>

            <v-col cols="6">
              <v-dialog v-model="timeEndDialog" width="290px">
                <template v-slot:activator="{ on, attrs }">
                  <v-text-field
                    v-model="bookingForm.time_end"
                    label="Check-out time"
                    append-icon="mdi-clock-time-four-outline"
                    readonly
                    v-bind="attrs"
                    v-on="on"
                  ></v-text-field>
                </template>
                <v-time-picker
                  v-if="timeEndDialog"
                  v-model="bookingForm.time_end"
                  full-width
                  color="#252525"
                  :min="
                    !bookingForm.date_range[1] ||
                    bookingForm.date_range[0] === bookingForm.date_range[1]
                      ? bookingForm.time_start
                      : `00:00`
                  "
                >
                </v-time-picker>
              </v-dialog>
            </v-col>
          </v-row>
        </v-card-actions>

        <v-divider></v-divider>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="#252525" class="white--text" @click="orderRoom">
            Book now
          </v-btn>
        </v-card-actions>

        <v-card-text class="green--text" v-if="bookingSuccess"
          >You have successfully placed your order! We have sent you an email
          about your order and A Meliã representative will contact you as soon
          as possible.</v-card-text
        >
        <v-card-text class="red--text" v-if="bookingFail"
          >The room was fully booked during that time</v-card-text
        >
      </v-card>
    </v-dialog>
  </v-app>
</template>

<script>
import numeral from "numeral";
import HomeBlogs from "../components/HomeBlogs.vue";

import { format } from "date-fns";

export default {
  components: {HomeBlogs },

  data: () => ({
    userInfor: {},
    isAdmin: false,
    isLogged: false,

    bookingDialog: false,
    date: new Date(new Date().getTime() + 24 * 60 * 60 * 1000),

    slideItems: [],

    timeStartDialog: false,
    timeEndDialog: false,
    dateRangeDialog: false,

    bookingForm: {
      select: {},
      name: "",
      phone_number: "",
      email: "",
      time_start: "",
      time_end: "",
      date_range: [],
    },

    bookingSuccess: false,
    bookingFail: false,
  }),

  computed: {},

  methods: {
    goToMyOrders() {
      this.$router.push({ name: "orders" });
    },
    
    logOut() {
      localStorage.clear();
      window.location.reload();
    },

    openFacebookFanpage() {
      const urlToOpen =
        "https://www.facebook.com/profile.php?id=100012623027984";
      window.open(urlToOpen, "_blank");
    },

    openInstaFanpage() {
      const urlToOpen = "https://www.instagram.com/_tao_meo_/";
      window.open(urlToOpen, "_blank");
    },

    openLoginPage() {
      this.$router.push({ name: "log-in" });
    },

    formattedPrice(price) {
      return numeral(price).format("0,0");
    },

    getRoomTypesList() {
      this.$axios
        .get("/room_types/all")
        .then((res) => {
          this.slideItems = res.data;
        })
        .catch((err) => {
          console.error(err);
        });
    },

    showBookingForm(item) {
      this.bookingForm.select = item;
      this.bookingDialog = true;
    },

    orderRoom() {
      this.$axios
        .post(`/orders`, {
          name: this.bookingForm.name,
          phone_number: this.bookingForm.phone_number,
          email: this.bookingForm.email,
          time_start:
            this.bookingForm.date_range[0] + " " + this.bookingForm.time_start,
          time_end: !this.bookingForm.date_range[1]
            ? this.bookingForm.date_range[0]
            : this.bookingForm.date_range[1] + " " + this.bookingForm.time_end,
          room_type_id: this.bookingForm.select.id,
        })

        .then((res) => {
          this.bookingSuccess = true;
          setTimeout(() => {
            this.bookingDialog = false;
            this.bookingSuccess = false;
            this.bookingForm = {
              select: {},
              name: "",
              phone_number: "",
              time_start: "",
              time_end: "",
              date_range: [],
            };
          }, 3000);
        })

        .catch((err) => {
          this.bookingFail = true;
          setTimeout(() => {
            this.bookingDialog = false;
            this.bookingFail = false;
            this.bookingForm = {
              select: {},
              name: "",
              phone_number: "",
              time_start: "",
              time_end: "",
              date_range: [],
            };
          }, 3000);
        });
    },
  },

  created() {
    this.getRoomTypesList();
    this.date = format(this.date, "yyyy-MM-dd");
    this.userInfor = this.$store.state.userInfor;
    this.isAdmin = this.$store.state.isAdmin;
    this.isLogged = this.$store.state.isLogged;

    this.bookingForm = {
      select: {},
      name: this.userInfor.name || "",
      phone_number: this.userInfor.phone_number || "",
      email: this.userInfor.email || "",
      time_start: "",
      time_end: "",
      date_range: [],
    };
  },
};
</script>

<style lang="scss" scoped>
.remove-padding {
  padding: 0 !important;
}

.icon-hover {
  color: #252525;

  &:hover {
    opacity: 0.5;
  }
}
</style>
