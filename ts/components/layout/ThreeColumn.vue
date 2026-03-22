<template>
	<div class="h-100">
		<div class="container-fluid h-100">
			<div class="row h-100">
				<div class="col-md-3 p-0">
					<slot name="left"></slot>
				</div>
				<div class="col-md-6 p-0">
					<slot name="center"></slot>
				</div>
				<div class="col-md-3 p-0">
					<slot name="right"></slot>
				</div>
			</div>
		</div>
	</div>
</template>
