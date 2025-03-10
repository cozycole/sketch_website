import { VideoCarousel } from '../components/videoCarousel.js'

export function initHome() {
  const carouselSections = document.getElementsByClassName('carouselSection');
  for (let section of carouselSections) {
    new VideoCarousel(
      section.querySelector('.carousel'),
      section.querySelector('.carouselPrevBtn'),
      section.querySelector('.carouselNextBtn'),
    );
  }
}
