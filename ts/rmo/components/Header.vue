<style scoped>
.banner {
	display: block;
}
.banner-container {
	position: relative;
	width: 100%;
	background-color: #f76436;
	overflow: hidden;
}
@include media-breakpoint-up(xs) {
	.banner {
		width: 100%;
	}
}
@include media-breakpoint-up(xl) {
	.banner {
		height: 100%;
	}
}
</style>
<template>
	<section class="py-2 bg-primary text-white">
		<div class="banner-container d-flex justify-content-center">
			<img class="banner" src="/static/img/rmo/banner.jpg" />
		</div>
	</section>
</template>
