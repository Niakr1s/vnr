import { Component, OnDestroy, OnInit } from '@angular/core';
import { Subscription } from 'rxjs';
import { Sentence } from '../services/models/sentence';
import { SentenceService } from '../services/sentence.service';

@Component({
  selector: 'app-sentence',
  templateUrl: './sentence.component.html',
  styleUrls: ['./sentence.component.css'],
})
export class SentenceComponent implements OnInit, OnDestroy {
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
