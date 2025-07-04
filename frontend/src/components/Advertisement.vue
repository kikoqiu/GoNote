<!-- @format -->

<template>
  <div class="advertisement-moudle">
    <el-carousel :height="carouselHeight" type="card" :interval="5000" arrow="always">
      <el-carousel-item v-for="item in items" :key="item.alt">
        <a :href="item.path" target="_blank" rel="noreferrer noopener">
          <img class="advertisement-img" :src="item.image" :alt="item.alt" />
        </a>
      </el-carousel-item>
    </el-carousel>
  </div>
</template>

<script>
export default {
  name: 'Advertisement',

  data() {
    return {
      carouselHeight: '100px',
      items: [
      ]
    }
  },

  computed: {},

  components: {},

  mounted() {},

  methods: {}
}
</script>

<style lang="less">
@import './../assets/styles/style.less';
.advertisement-moudle {
  margin-bottom: 6 * @size-factor;
  .text-center {
    text-align: center;
  }
  .advertisement-img {
    width: 100%;
    min-height: 60px;
  }
  .VueCarousel {
    .VueCarousel-pagination {
      .VueCarousel-dot {
        .VueCarousel-dot-button {
          width: 12px !important;
          height: 12px !important;
          &:focus {
            border: none !important;
            outline: none !important;
          }
        }
      }
    }
  }
}
</style>
