import { Component, OnInit } from '@angular/core';
import { Subscription } from 'rxjs';
import { Sentence } from '../services/models/sentence';
import { SentenceService } from '../services/sentence.service';

@Component({
  selector: 'app-translation',
  templateUrl: './translation.component.html',
  styleUrls: ['./translation.component.css'],
})
export class TranslationComponent implements OnInit {
  sentence!: Sentence | null;
  private sentenceSub!: Subscription;

  constructor(private sentenceService: SentenceService) {}
  ngOnInit(): void {
    this.sentenceSub = this.sentenceService.currentSentence$.subscribe({
      next: (sentence) => {
        this.sentence = sentence;
      },
    });
  }

  ngOnDestroy(): void {
    this.sentenceSub.unsubscribe();
  }
}
